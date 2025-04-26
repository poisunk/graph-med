package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"graph-med/app/captcha/rpc/captcha"

	"graph-med/app/captcha/api/internal/svc"
	"graph-med/app/captcha/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// verify captcha
func NewVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyLogic {
	return &VerifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyLogic) Verify(req *types.VerifyCaptchaReq) (resp *types.VerifyCaptchaResp, err error) {
	verifyResp, err := l.svcCtx.CaptchaRpc.VerifyCaptcha(l.ctx, &captcha.VerifyCaptchaReq{
		CaptchaType: req.CaptchaType,
		CaptchaId:   req.CaptchaId,
		Answer:      req.Answer,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.VerifyCaptchaResp{}
	_ = copier.Copy(resp, verifyResp)

	return resp, nil
}
