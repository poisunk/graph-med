package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf

	EmailConfig struct {
		Host     string
		Port     int
		Username string
		Password string
	}
}
