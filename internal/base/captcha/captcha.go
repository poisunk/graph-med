package captcha

import "errors"

const (
	DigitType  = "digit"
	RotateType = "rotate"
	EmailType  = "email"
)

func GetCaptchaData(captchaType string, params ...interface{}) (interface{}, error) {
	switch captchaType {
	case DigitType:
		return getDigitCaptchaData()
	case RotateType:
		return getRotateCaptchaData()
	case EmailType:
		if len(params) == 0 {
			return nil, errors.New("email is required")
		}
		return getEmailCaptchaData(params[0].(string))
	default:
		return DigitCaptchaData{}, nil
	}
}

func VerifyCaptcha(captchaType, id, answer string) bool {
	switch captchaType {
	case DigitType:
		return verifyDigitCaptcha(id, answer)
	case RotateType:
		return verifyRotateCaptchaData(id, answer)
	case EmailType:
		return verifyEmailCaptcha(id, answer)
	default:
		return false
	}
}
