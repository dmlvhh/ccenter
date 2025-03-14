package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Mysql struct {
		DataSource string
	}
	Etcd struct {
		Hosts    []string
		Key      string
		LeaseTTL int64
	}
}
