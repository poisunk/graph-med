package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"graph-med/app/captcha/api/internal/svc"
	"graph-med/app/captcha/api/internal/types"
	"graph-med/app/captcha/rpc/captcha"

	"github.com/zeromicro/go-zero/core/logx"
)

type DigitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// digit captcha
func NewDigitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DigitLogic {
	return &DigitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DigitLogic) Digit(req *types.DigitCaptchaReq) (resp *types.DigitCaptchaResp, err error) {
	digitResp, err := l.svcCtx.CaptchaRpc.DigitCaptcha(l.ctx, &captcha.DigitCaptchaReq{})
	if err != nil {
		return nil, err
	}

	resp = &types.DigitCaptchaResp{}
	_ = copier.Copy(resp, digitResp)

	return resp, nil
}
