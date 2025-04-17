package service

import (
	"errors"
	"github.com/casbin/casbin/v2"
	"graph-med/internal/base/captcha"
	casbinRole "graph-med/internal/base/casbin"
	"graph-med/internal/base/redis"
	"graph-med/internal/model"
	"graph-med/internal/repository"
	"graph-med/internal/utils"
	"time"
)

type AuthService struct {
	userRepo *repository.UserRepository
	enf      *casbin.Enforcer
}

func NewAuthService(userRepo *repository.UserRepository, enf *casbin.Enforcer) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		enf:      enf,
	}
}

// RegisterUser 注册用户
func (a *AuthService) RegisterUser(username, email, password string) (string, error) {
	_, err := a.userRepo.FindByEmail(email)
	if err == nil {
		return "", errors.New("user already exists")
	}

	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	userID := utils.GenerateUUID()

	user := &model.ChatUser{
		UserID:   userID,
		Username: username,
		Email:    email,
		Password: hashPassword,
	}
	err = a.userRepo.CreateUser(user)
	if err != nil {
		return "", err
	}

	token, err := utils.GenerateToken(userID, username)
	if err != nil {
		return "", err
	}

	err = redis.Set("user-"+token, user, time.Hour)
	if err != nil {
		return "", err
	}

	// 初始化用户策略
	err = a.InitializeUserPolicy(user.UserID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// LoginUser 登录用户
func (a *AuthService) LoginUser(email, password string) (string, error) {
	user, err := a.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("用户不存在")
	}

	if !utils.ComparePasswords(user.Password, password) {
		return "", errors.New("用户名或密码错误")
	}

	token, err := utils.GenerateToken(user.UserID, user.Username)
	if err != nil {
		return "", err
	}

	err = redis.Set("user-"+token, user, time.Hour)
	if err != nil {
		return "", err
	}

	return token, nil
}

// InitializeUserPolicy 初始化用户策略
func (a *AuthService) InitializeUserPolicy(userID string) error {
	suc, err := a.enf.AddGroupingPolicy(userID, casbinRole.NORMAL_USER)
	if err != nil {
		return err
	}

	if !suc {
		return errors.New("failed to add policy")
	}

	return nil
}

// GetCaptchaData 获取验证码
func (a *AuthService) GetCaptchaData(captchaType string, params ...interface{}) (interface{}, error) {
	return captcha.GetCaptchaData(captchaType, params...)
}

// ValidateCaptcha 验证验证码
func (a *AuthService) ValidateCaptcha(captchaType, id, answer string) bool {
	return captcha.VerifyCaptcha(captchaType, id, answer)
}

// LogoutUser 登出用户
func (a *AuthService) LogoutUser(token string) error {
	err := redis.Del("user-" + token)
	if err != nil {
		return err
	}
	return nil
}
