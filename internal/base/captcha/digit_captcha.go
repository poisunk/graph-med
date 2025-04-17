package captcha

import (
	"graph-med/internal/base/redis"
	"time"

	"github.com/mojocn/base64Captcha"
)

type DigitCaptchaData struct {
	CaptchaId string `json:"captcha_id"`       //验证码id
	Data      string `json:"data"`             //验证码数据base64类型
	Answer    string `json:"answer,omitempty"` //验证码数字
}

type RedisStore struct{}

func (s RedisStore) Set(id string, value string) error {
	err := redis.Set("captcha-"+id, value, 3*time.Minute)
	return err
}

func (s RedisStore) Get(id string, clear bool) string {
	var value string
	value, err := redis.Get[string]("captcha-" + id)
	if clear {
		err = redis.Del(id)
		if err != nil {
			return ""
		}
	}
	return value
}

func (s RedisStore) Verify(id, answer string, clear bool) bool {
	value := s.Get(id, clear)
	return value == answer
}

// 数字驱动
var digitDriver = base64Captcha.DriverDigit{
	Height:   50,
	Width:    200,
	Length:   4,   //验证码长度
	MaxSkew:  0.7, //倾斜
	DotCount: 1,   //背景的点数，越大，字体越模糊
}

// redis Store
var redisStore = RedisStore{}

func getDigitCaptchaData() (DigitCaptchaData, error) {
	var ret DigitCaptchaData
	c := base64Captcha.NewCaptcha(&digitDriver, redisStore)
	id, b64s, answer, err := c.Generate()
	if err != nil {
		return ret, err
	}
	ret.CaptchaId = id
	ret.Data = b64s

	err = redisStore.Set(id, answer)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func verifyDigitCaptcha(id, answer string) bool {
	return redisStore.Verify(id, answer, true)
}
