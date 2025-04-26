package logic

import (
	"context"
	"fmt"
	"github.com/mojocn/base64Captcha"
	"github.com/zeromicro/go-zero/core/logx"
	"graph-med/app/captcha/rpc/internal/svc"
	"graph-med/app/captcha/rpc/pd"
)

type DigitCaptchaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDigitCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DigitCaptchaLogic {
	return &DigitCaptchaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DigitCaptchaLogic) DigitCaptcha(in *pd.DigitCaptchaReq) (*pd.DigitCaptchaResp, error) {
	c := base64Captcha.NewCaptcha(&base64Captcha.DriverDigit{
		Height:   50,
		Width:    200,
		Length:   4,   // 验证码长度
		MaxSkew:  0.7, // 倾斜
		DotCount: 1,   // 背景的点数，越大，字体越模糊
	}, base64Captcha.DefaultMemStore)

	id, b64s, answer, err := c.Generate()
	if err != nil {
		return nil, err
	}

	graphMedCaptchaKey := fmt.Sprintf("%s%v", cacheGraphMedCaptchaDigitPrefix, id)
	err = l.svcCtx.RedisClient.Setex(graphMedCaptchaKey, answer, captchaExpire)
	if err != nil {
		return nil, err
	}

	return &pd.DigitCaptchaResp{
		CaptchaId:  id,
		CaptchaImg: b64s,
	}, nil
}

func (l *DigitCaptchaLogic) VerifyDigitCaptcha(id, answer string) (bool, error) {
	graphMedCaptchaKey := fmt.Sprintf("%s%v", cacheGraphMedCaptchaDigitPrefix, id)
	value, err := l.svcCtx.RedisClient.Get(graphMedCaptchaKey)
	if err != nil {
		return false, err
	}

	return value == answer, nil
}
