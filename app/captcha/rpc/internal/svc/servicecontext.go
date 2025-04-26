package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"graph-med/app/captcha/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config

	RedisClient *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		RedisClient: redis.MustNewRedis(c.Redis.RedisConf),
	}
}
