package etcd

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/netx"
	"go-zero2/initialize"
	"strings"
)

// DeliveryAddress 将服务地址注册到 etcd 中。
// etcdAddr: etcd 服务的地址。
// serviceName: 要注册的服务名称。
// addr: 要注册的服务地址。
func DeliveryAddress(etcdAddr string, serviceName string, addr string) {
	// 根据冒号分隔服务地址
	list := strings.Split(addr, ":")
	// 检查地址格式是否正确
	if len(list) != 2 {
		logx.Errorf("invalid addr: %s", addr)
		return
	}

	// 如果服务地址是 "0.0.0.0"，则替换为本机 IP
	if list[0] == "0.0.0.0" {
		ip := netx.InternalIp() // 获取本机 IP
		strings.ReplaceAll(addr, "0.0.0.0", ip)
	}
	// 初始化 etcd 客户端
	client := initialize.InitEtcd(etcdAddr)
	// 将服务地址写入 etcd
	_, err := client.Put(context.Background(), serviceName, addr)
	if err != nil {
		logx.Errorf("put etcd err: %s", err.Error())
		return
	}
	logx.Infof("put etcd success: %s %s", serviceName, addr)
}

// GetServiceAddr 从 etcd 中获取服务地址。
// etcdAddr: etcd 服务的地址。
// serviceName: 要查询的服务名称。
// 返回值: 服务的地址字符串。如果无法获取到服务地址，返回空字符串。
func GetServiceAddr(etcdAddr string, serviceName string) (addr string, err error) {
	// 初始化 etcd 客户端
	client := initialize.InitEtcd(etcdAddr)
	// 从 etcd 中获取服务地址
	resp, err := client.Get(context.Background(), serviceName)
	// 处理获取服务地址的结果
	if err == nil && len(resp.Kvs) > 0 {
		return string(resp.Kvs[0].Value), nil
	}
	return
}
