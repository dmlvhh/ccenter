package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"

	"go-zero2/sub2/user_api/internal/config"
	"go-zero2/sub2/user_api/internal/handler"
	"go-zero2/sub2/user_api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/userapi-api.yaml", "the config file")

func registerToEtcd(c config.Config) error {
	// 创建 etcd 客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.Etcd.Hosts,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("连接 etcd 失败: %v", err)
	}
	defer cli.Close()
	// 注册服务地址
	serviceAddr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	_, err = cli.Put(context.Background(), c.Etcd.Key, serviceAddr)
	if err != nil {
		return fmt.Errorf("注册服务失败: %v", err)
	}

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
