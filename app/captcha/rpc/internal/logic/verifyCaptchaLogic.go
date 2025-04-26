package logic

import (
	"context"
	"fmt"
	"github.com/wenlng/go-captcha/v2/rotate"
	"strconv"

	"graph-med/app/captcha/rpc/internal/svc"
	"graph-med/app/captcha/rpc/pd"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	cacheGraphMedCaptchaDigitPrefix  = "cache:graphMedCaptcha:digit:"
	cacheGraphMedCaptchaEmailPrefix  = "cache:graphMedCaptcha:email:"
	cacheGraphMedCaptchaRotatePrefix = "cache:graphMedCaptcha:rotate:"

	captchaExpire = 300
)

type VerifyCaptchaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyCaptchaLogic {
	return &VerifyCaptchaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyCaptchaLogic) VerifyCaptcha(in *pd.VerifyCaptchaReq) (*pd.VerifyCaptchaResp, error) {
	var captchaKeyPrefix string
	switch in.CaptchaType {
	case "digit":
		captchaKeyPrefix = cacheGraphMedCaptchaDigitPrefix
		break
	case "email":
		captchaKeyPrefix = cacheGraphMedCaptchaEmailPrefix
		break
	case "rotate":
		captchaKeyPrefix = cacheGraphMedCaptchaRotatePrefix
		break
	default:
		return nil, fmt.Errorf("invalid captcha type: %s", in.CaptchaType)
	}

	captchaKey := fmt.Sprintf("%s%v", captchaKeyPrefix, in.CaptchaId)

	// 从Redis中获取验证码
	value, err := l.svcCtx.RedisClient.Get(captchaKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get redis: %v", err)
	}

	var success bool
	switch in.CaptchaType {
	case "digit":
	case "email":
		success = value == in.Answer
		break
	case "rotate":
		srcAngle, err := l.svcCtx.RedisClient.Get(captchaKey)
		if err != nil {
			return nil, err
		}

		srcAngleInt, err := strconv.ParseInt(srcAngle, 10, 64)
		if err != nil {
			return nil, err
		}
		angle, err := strconv.ParseInt(in.Answer, 10, 64)
		if err != nil {
			return nil, err
		}

		success = rotate.CheckAngle(srcAngleInt, angle, 10)
		break
	default:
		return nil, fmt.Errorf("invalid captcha type: %s", in.CaptchaType)
	}
	if err != nil {
		return nil, err
	}

	// 验证成功后删除Redis中的验证码，防止重复使用
	_, err = l.svcCtx.RedisClient.Del(captchaKey)
	if err != nil {
		return nil, fmt.Errorf("failed to del redis: %v", err)
	}

	return &pd.VerifyCaptchaResp{
		Success: success,
	}, nil
}
