package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"graph-med/app/usercenter/model"
	"graph-med/app/usercenter/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config

	RedisClient *redis.Redis

	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)

	return &ServiceContext{
		Config:      c,
		RedisClient: redis.MustNewRedis(c.Redis.RedisConf),
		UserModel:   model.NewUserModel(sqlConn, c.Cache),
	}
}
