package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"graph-med/app/captcha/rpc/captcha"

	"graph-med/app/captcha/api/internal/svc"
	"graph-med/app/captcha/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RotateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// rotate captcha
func NewRotateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RotateLogic {
	return &RotateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RotateLogic) Rotate(req *types.RotateCaptchaReq) (resp *types.RotateCaptchaResp, err error) {
	rotateResp, err := l.svcCtx.CaptchaRpc.RotateCaptcha(l.ctx, &captcha.RotateCaptchaReq{})
	if err != nil {
		return nil, err
	}

	resp = &types.RotateCaptchaResp{}
	_ = copier.Copy(resp, rotateResp)

	return resp, nil
}
