package captcha

import (
	"fmt"
	images "github.com/wenlng/go-captcha-assets/resources/images_v2"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/rotate"
	"graph-med/internal/base/redis"
	"graph-med/internal/utils"
	"log"
	"strconv"
	"time"
)

var rotateCapt rotate.Captcha

type RotateCaptchaData struct {
	CaptchaId   string `json:"captcha_id"`
	ImageBase64 string `json:"image_base64"`
	ThumbBase64 string `json:"thumb_base64"`
	ParentSize  int    `json:"parent_size"`
	ChildSize   int    `json:"child_size"`
}

func init() {
	// 初始化 rotate captcha
	builder := rotate.NewBuilder(rotate.WithRangeAnglePos([]option.RangeVal{
		{Min: 20, Max: 330},
	}))

	imgs, err := images.GetImages()
	if err != nil {
		log.Fatalln(err)
	}

	builder.SetResources(
		rotate.WithImages(imgs),
	)

	rotateCapt = builder.Make()
}

func getRotateCaptchaData() (*RotateCaptchaData, error) {
	captData, err := rotateCapt.Generate()
	if err != nil {
		return nil, err
	}

	captchaKey := utils.GenerateUUID()
	angle := captData.GetData().Angle

	err = redis.Set("captcha-"+captchaKey, angle, 3*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to set redis: %v", err)
	}

	imageBase64, err := captData.GetMasterImage().ToBase64()
	if err != nil {
		return nil, err
	}

	thumbBase64, err := captData.GetThumbImage().ToBase64()
	if err != nil {
		return nil, err
	}

	data := &RotateCaptchaData{
		CaptchaId:   captchaKey,
		ImageBase64: imageBase64,
		ThumbBase64: thumbBase64,
		ParentSize:  captData.GetData().ParentWidth,
		ChildSize:   captData.GetData().Width,
	}
	return data, nil
}

func verifyRotateCaptchaData(id, answer string) bool {
	srcAngle, err := redis.Get[int64]("captcha-" + id)
	if err != nil {
		return false
	}
	err = redis.Del("captcha-" + id)
	if err != nil {
		return false
	}

	angle, err := strconv.ParseInt(answer, 10, 64)
	if err != nil {
		return false
	}
	return rotate.CheckAngle(srcAngle, angle, 10)
}
