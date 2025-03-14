package initialize

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// InitEtcd 初始化并返回一个Etcd客户端实例。
func InitEtcd(add string) *clientv3.Client {
	// 创建一个Etcd客户端实例
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{add},   // 指定Etcd服务器地址
		DialTimeout: 5 * time.Second, // 设置连接超时时间
	})
	if err != nil {
		panic(err) // 如果连接创建失败，则抛出异常
	}
	return cli
}
