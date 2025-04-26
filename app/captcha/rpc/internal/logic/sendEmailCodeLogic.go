package logic

import (
	"context"
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"

	"graph-med/app/captcha/rpc/internal/svc"
	"graph-med/app/captcha/rpc/pd"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendEmailCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendEmailCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendEmailCodeLogic {
	return &SendEmailCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendEmailCodeLogic) SendEmailCode(in *pd.SendEmailCodeReq) (*pd.SendEmailCodeResp, error) {
	if in.Email == "" {
		return nil, fmt.Errorf("邮箱地址为空")
	}

	smtpServer := l.svcCtx.Config.EmailConfig.Host
	smtpPort := l.svcCtx.Config.EmailConfig.Port
	smtpAddr := fmt.Sprintf("%s:%d", smtpServer, smtpPort)

	hostEmail := l.svcCtx.Config.EmailConfig.Username
	hostPass := l.svcCtx.Config.EmailConfig.Password

	// 创建邮件对象
	e := email.NewEmail()
	e.From = fmt.Sprintf("GraphMed <%s>", hostEmail)
	e.To = []string{in.Email}
	e.Subject = "验证码"
	e.HTML = []byte(fmt.Sprintf(`
		<div style="max-width: 600px; margin: 0 auto; padding: 20px; font-family: Arial, sans-serif;">
			<h2>验证码</h2>
			<p>您的验证码是：<strong style="font-size: 24px; color: #1a73e8;">%s</strong></p>
			<p>验证码有效期为3分钟，请尽快完成验证。</p>
			<p style="color: #666;">如果这不是您的操作，请忽略此邮件。</p>
		</div>
	`, in.Code))

	ret := &pd.SendEmailCodeResp{
		Success: true,
	}

	// 配置SMTP服务器信息
	err := e.Send(smtpAddr, smtp.PlainAuth("",
		hostEmail,
		hostPass,
		smtpServer,
	))
	if err != nil {
		ret.Success = false
	}

	return ret, nil
}
