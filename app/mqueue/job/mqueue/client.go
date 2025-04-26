package mqueue

import (
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func NewAsynqClient(c redis.RedisConf) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: c.Host, Password: c.Pass})
}
