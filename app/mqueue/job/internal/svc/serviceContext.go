package svc

import (
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/zrpc"
	"graph-med/app/captcha/rpc/captcha"
	"graph-med/app/mqueue/job/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	AsynqServer *asynq.Server

	CaptchaRpc captcha.Captcha
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		AsynqServer: newAsynqServer(c),

		CaptchaRpc: captcha.NewCaptcha(zrpc.MustNewClient(c.CaptchaRpcConf)),
	}
}
