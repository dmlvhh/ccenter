package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"net"
	"time"

	"go-zero2/sub1/user_api/internal/config"
	"go-zero2/sub1/user_api/internal/handler"
	"go-zero2/sub1/user_api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/userapi-api.yaml", "the config file")

// getLocalIP 获取本机 IP 地址
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("无法获取本机 IP 地址")
}
func registerToEtcd(c config.Config) error {
	ip, err := getLocalIP()
	// 服务地址
	serviceAddr := fmt.Sprintf("%s:%d", ip, c.Port)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.Etcd.Hosts,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("连接 etcd 失败: %v", err)
	}
	defer cli.Close()

	// 创建租约
	leaseResp, err := cli.Grant(context.Background(), c.Etcd.LeaseTTL)
	if err != nil {
		return fmt.Errorf("创建租约失败: %v", err)
	}

	// 注册服务地址
	//serviceAddr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	_, err = cli.Put(context.Background(), c.Etcd.Key, serviceAddr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return fmt.Errorf("注册服务失败: %v", err)
	}

	// 定期续期
	go func() {
		for {
			_, err2 := cli.KeepAliveOnce(context.Background(), leaseResp.ID)
			if err2 != nil {
				logx.Errorf("续期失败: %v", err2)
				return
			}
			time.Sleep(time.Duration(c.Etcd.LeaseTTL/2) * time.Second)
		}
	}()

	logx.Infof("服务注册成功: %s -> %s", c.Etcd.Key, serviceAddr)
	return nil
}

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	// 注册服务到 etcd
	if err := registerToEtcd(c); err != nil {
		logx.Error("注册服务失败:", err)
		return
	}

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
