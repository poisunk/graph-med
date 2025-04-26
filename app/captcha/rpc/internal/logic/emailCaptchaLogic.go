package logic

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"graph-med/app/captcha/rpc/internal/svc"
	"graph-med/app/captcha/rpc/pd"
	"graph-med/internal/base/redis"
	"graph-med/pkg/tool"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmailCaptchaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEmailCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailCaptchaLogic {
	return &EmailCaptchaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EmailCaptchaLogic) EmailCaptcha(in *pd.EmailCaptchaReq) (*pd.EmailCaptchaResp, error) {
	captchaId := uuid.New().String()

	graphMedCaptchaEmailKey := fmt.Sprintf("%s%v", cacheGraphMedCaptchaEmailPrefix, captchaId)

	// 检查是否在3分钟内已经发送过验证码
	//value, err := l.svcCtx.RedisClient.Get(graphMedCaptchaEmailKey)
	//if err == nil && value != "" {
	//	return nil, fmt.Errorf("email already sent")
	//}

	// 生成6位随机验证码
	code := tool.RandomString(6, tool.RAND_STRING_KIND_NUM)

	// 存储验证码到Redis，有效期3分钟
	err := l.svcCtx.RedisClient.Setex(graphMedCaptchaEmailKey, code, captchaExpire)
	if err != nil {
		return nil, fmt.Errorf("failed to set redis: %v", err)
	}

	// 发送验证码到邮箱
	//err = sendEmailCode(email, code)
	if err != nil {
		// 如果发送失败，删除Redis中的验证码
		_, _ = l.svcCtx.RedisClient.Del(graphMedCaptchaEmailKey)
		return nil, fmt.Errorf("failed to send email: %v", err)
	}

	return &pd.EmailCaptchaResp{
		CaptchaId: captchaId,
	}, nil
}

func (l *EmailCaptchaLogic) VerifyCaptcha(id, answer string) (bool, error) {
	graphMedCaptchaKey := fmt.Sprintf("%s%v", cacheGraphMedCaptchaEmailPrefix, id)

	value, err := l.svcCtx.RedisClient.Get(graphMedCaptchaKey)
	if err != nil {
		return false, err
	}

	// 验证成功后删除Redis中的验证码，防止重复使用
	err = redis.Del(graphMedCaptchaKey)
	if err != nil {
		return false, err
	}

	return value == answer, nil
}
