package handler

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"graph-med/app/captcha/rpc/captcha"
	"graph-med/app/mqueue/job/internal/svc"
	"graph-med/app/mqueue/job/jobtype"
	"graph-med/pkg/xerr"
)

var ErrSendEmailFail = xerr.NewErrMsg("发送邮件失败")

type SendEmailHandler struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewSendEmailHandler(svcCtx *svc.ServiceContext) *SendEmailHandler {
	return &SendEmailHandler{
		svcCtx: svcCtx,
	}
}

func (s *SendEmailHandler) ProcessTask(cxt context.Context, task *asynq.Task) error {

	var p jobtype.CaptchaSendEmailPayload
	if err := json.Unmarshal(task.Payload(), &p); err != nil {
		return errors.Wrapf(ErrSendEmailFail, "send email json unmarshal err:%v, payLoad:%+v", err, task.Payload())
	}

	resp, err := s.svcCtx.CaptchaRpc.SendEmailCode(cxt, &captcha.SendEmailCodeReq{
		Email: p.Email,
		Code:  p.Code,
	})
	if err != nil {
		return errors.Wrapf(ErrSendEmailFail, "send email rpc err:%v", err)
	}

	if !resp.Success {
		s.Error(task, errors.Wrapf(ErrSendEmailFail, "send email rpc not success"))
	}

	return nil
}
