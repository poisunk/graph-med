package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"graph-med/app/captcha/api/internal/config"
	"graph-med/app/captcha/rpc/captcha"
)

type ServiceContext struct {
	Config config.Config

	CaptchaRpc captcha.Captcha
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		CaptchaRpc: captcha.NewCaptcha(zrpc.MustNewClient(c.CaptchaRpcConf)),
	}
}
