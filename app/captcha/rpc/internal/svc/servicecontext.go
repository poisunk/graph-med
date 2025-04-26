package svc

import (
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"graph-med/app/captcha/rpc/internal/config"
	"graph-med/app/mqueue/job/mqueue"
)

type ServiceContext struct {
	Config config.Config

	RedisClient *redis.Redis
	AsynqClient *asynq.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		RedisClient: redis.MustNewRedis(c.Redis.RedisConf),
		AsynqClient: mqueue.NewAsynqClient(c.Redis.RedisConf),
	}
}
