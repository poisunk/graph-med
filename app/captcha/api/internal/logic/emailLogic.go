package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"graph-med/app/captcha/rpc/captcha"

	"graph-med/app/captcha/api/internal/svc"
	"graph-med/app/captcha/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// email captcha
func NewEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailLogic {
	return &EmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmailLogic) Email(req *types.EmailCaptchaReq) (resp *types.EmailCaptchaResp, err error) {
	emailResp, err := l.svcCtx.CaptchaRpc.EmailCaptcha(l.ctx, &captcha.EmailCaptchaReq{
		Email: req.Email,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.EmailCaptchaResp{}
	_ = copier.Copy(resp, emailResp)

	return resp, nil
}
