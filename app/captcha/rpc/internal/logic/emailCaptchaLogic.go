package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"graph-med/app/captcha/rpc/internal/svc"
	"graph-med/app/captcha/rpc/pd"
	"graph-med/app/mqueue/job/jobtype"
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
	payload, err := json.Marshal(jobtype.CaptchaSendEmailPayload{
		Email: in.Email,
		Code:  code,
	})
	_, err = l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(jobtype.CaptchaSendEmail, payload))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("send email insert queue fail err :%+v , email : %s, code : %s", err, in.Email, code)

		// 如果发送失败，删除Redis中的验证码
		_, _ = l.svcCtx.RedisClient.Del(graphMedCaptchaEmailKey)
		return nil, fmt.Errorf("failed to send email: %v", err)
	}

	return &pd.EmailCaptchaResp{
		CaptchaId: captchaId,
	}, nil
}
