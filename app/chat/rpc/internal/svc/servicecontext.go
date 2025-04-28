package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"graph-med/app/chat/model"
	"graph-med/app/chat/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config

	RedisClient *redis.Redis

	ChatSessionModel  model.ChatSessionModel
	ChatMessageModel  model.ChatMessageModel
	ChatTypeModel     model.ChatTypeModel
	ChatFeedbackModel model.ChatFeedbackModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		RedisClient: redis.MustNewRedis(c.Redis.RedisConf),

		ChatSessionModel:  model.NewChatSessionModel(c.Mongo.Url, c.Mongo.DB, "ChatSession"),
		ChatMessageModel:  model.NewChatMessageModel(c.Mongo.Url, c.Mongo.DB, "ChatMessage"),
		ChatTypeModel:     model.NewChatTypeModel(c.Mongo.Url, c.Mongo.DB, "ChatType", c.Cache),
		ChatFeedbackModel: model.NewChatFeedbackModel(c.Mongo.Url, c.Mongo.DB, "ChatFeedback"),
	}
}
