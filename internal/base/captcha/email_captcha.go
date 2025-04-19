package captcha

import (
	"fmt"
	"graph-med/internal/base/logger"
	"graph-med/internal/base/redis"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"

	"github.com/jordan-wright/email"
	"go.uber.org/zap"
)

type EmailCaptchaData struct {
	CaptchaId string `json:"captcha_id"`     // 验证码ID
	Email     string `json:"email"`          // 邮箱地址
	Code      string `json:"code,omitempty"` // 验证码（仅在测试环境下返回）
}

var (
	initialized = false
	userEmail   string
	userPass    string
	smtpServer  string
	smtpPort    string
)

// SetupEmailCaptcha 设置邮件验证码
func SetupEmailCaptcha(email, password, host, port string) {
	initialized = true
	userEmail = email
	userPass = password
	smtpServer = host
	smtpPort = port
}

// generateRandomCode 生成随机验证码
func generateRandomCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < length; i++ {
		code += strconv.Itoa(rand.Intn(10))
	}
	return code
}

// sendEmailCode 发送邮件验证码
func sendEmailCode(emailAddr, code string) error {
	if !initialized {
		return fmt.Errorf("邮件验证码未初始化")
	}

	if emailAddr == "" {
		return fmt.Errorf("邮箱地址为空")
	}

	// 创建邮件对象
	e := email.NewEmail()
	e.From = fmt.Sprintf("GraphMed <%s>", userEmail)
	e.To = []string{emailAddr}
	e.Subject = "验证码"
	e.HTML = []byte(fmt.Sprintf(`
		<div style="max-width: 600px; margin: 0 auto; padding: 20px; font-family: Arial, sans-serif;">
			<h2>验证码</h2>
			<p>您的验证码是：<strong style="font-size: 24px; color: #1a73e8;">%s</strong></p>
			<p>验证码有效期为3分钟，请尽快完成验证。</p>
			<p style="color: #666;">如果这不是您的操作，请忽略此邮件。</p>
		</div>
	`, code))

	// 配置SMTP服务器信息

	go func() {
		err := e.Send(smtpServer+":"+smtpPort, smtp.PlainAuth("",
			userEmail,
			userPass,
			smtpServer,
		))

		if err != nil {
			logger.Error("邮件发送失败", zap.String("error", err.Error()))
			return
		}

		logger.Info("邮件发送成功")
	}()

	return nil
}

// 获取邮件验证码数据
func getEmailCaptchaData(email string) (interface{}, error) {
	// 检查是否在3分钟内已经发送过验证码
	value, err := redis.Get[string]("email-captcha-" + email)
	if err == nil && value != "" {
		return nil, fmt.Errorf("email already sent")
	}

	// 生成6位随机验证码
	code := generateRandomCode(6)

	// 存储验证码到Redis，有效期3分钟
	err = redis.Set("email-captcha-"+email, code, 3*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to set redis: %v", err)
	}

	// 发送验证码到邮箱
	err = sendEmailCode(email, code)
	if err != nil {
		// 如果发送失败，删除Redis中的验证码
		_ = redis.Del("email-captcha-" + email)
		return nil, fmt.Errorf("failed to send email: %v", err)
	}

	return nil, nil
}

// 验证邮件验证码
func verifyEmailCaptcha(email, answer string) bool {
	value, err := redis.Get[string]("email-captcha-" + email)
	if err != nil {
		return false
	}

	// 验证成功后删除Redis中的验证码，防止重复使用
	err = redis.Del("email-captcha-" + email)
	if err != nil {
		return false
	}

	return value == answer
}
