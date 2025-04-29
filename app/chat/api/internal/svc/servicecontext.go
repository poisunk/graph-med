package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"graph-med/app/chat/api/internal/config"
	"graph-med/app/chat/rpc/chat"
)

type ServiceContext struct {
	Config config.Config

	ChatRpc chat.Chat
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		ChatRpc: chat.NewChat(zrpc.MustNewClient(c.ChatRpcConf)),
	}
}
