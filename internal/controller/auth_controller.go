package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"graph-med/internal/base/captcha"
	"graph-med/internal/base/response"
	"graph-med/internal/service"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Login 处理登录请求
func (c *AuthController) Login(ctx *gin.Context) {
	var loginReq struct {
		Email         string `json:"email" binding:"required"`
		Password      string `json:"password" binding:"required"`
		CaptchaType   string `json:"captcha_type" binding:"required"`
		CaptchaID     string `json:"captcha_id" binding:"required"`
		CaptchaAnswer string `json:"captcha_answer" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		response.Error(ctx, err)
		return
	}

	valid := c.authService.ValidateCaptcha(loginReq.CaptchaType, loginReq.CaptchaID, loginReq.CaptchaAnswer)
	if !valid {
		response.Error(ctx, errors.New("invalid captcha"))
		return
	}

	token, err := c.authService.LoginUser(loginReq.Email, loginReq.Password)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, gin.H{"token": token})
}

// Register 处理注册请求
func (c *AuthController) Register(ctx *gin.Context) {
	var registerReq struct {
		Username      string `json:"username" binding:"required"`
		Email         string `json:"email" binding:"required"`
		Password      string `json:"password" binding:"required"`
		CaptchaType   string `json:"captcha_type" binding:"required"`
		CaptchaID     string `json:"captcha_id" binding:"required"`
		CaptchaAnswer string `json:"captcha_answer" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&registerReq); err != nil {
		response.Error(ctx, err)
		return
	}

	valid := c.authService.ValidateCaptcha(registerReq.CaptchaType, registerReq.CaptchaID, registerReq.CaptchaAnswer)
	if !valid {
		response.Error(ctx, errors.New("invalid captcha"))
		return
	}

	toekn, err := c.authService.RegisterUser(registerReq.Username, registerReq.Email, registerReq.Password)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, gin.H{"token": toekn})
}

// GetCaptcha 获取验证码
func (a *AuthController) GetCaptcha(ctx *gin.Context) {
	captchaType := ctx.DefaultQuery("type", captcha.DigitType)

	captResp, err := a.authService.GetCaptchaData(captchaType)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, captResp)
}

// GetEmailCaptcha 获取邮箱验证码
func (a *AuthController) GetEmailCaptcha(ctx *gin.Context) {
	var captReq struct {
		Email string `json:"email" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&captReq); err != nil {
		response.Error(ctx, err)
		return
	}

	captResp, err := a.authService.GetCaptchaData(captcha.EmailType, captReq.Email)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, captResp)
}

// ValidateCaptcha 验证验证码
func (a *AuthController) ValidateCaptcha(ctx *gin.Context) {
	var captReq struct {
		Type       string `json:"type"`
		CaptchaKey string `json:"captcha_key" binding:"required"`
		Answer     string `json:"answer" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&captReq); err != nil {
		response.Error(ctx, err)
		return
	}

	if captReq.Type == "" {
		captReq.Type = captcha.DigitType
	}

	valid := a.authService.ValidateCaptcha(captReq.Type, captReq.CaptchaKey, captReq.Answer)
	if !valid {
		response.Error(ctx, errors.New("invalid captcha"))
		return
	}

	response.Success(ctx, nil)
}

// Logout 处理登出请求
func (a *AuthController) Logout(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Access-Token")
	if token == "" {
		response.Unauthorized(ctx)
		return
	}

	err := a.authService.LogoutUser(token)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}
