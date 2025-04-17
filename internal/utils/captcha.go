package utils

import (
	images "github.com/wenlng/go-captcha-assets/resources/images_v2"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/rotate"
	"log"
)

var rotateCapt rotate.Captcha

func init() {
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

func GetRotateCaptchaData() (rotate.CaptchaData, error) {
	captData, err := rotateCapt.Generate()
	if err != nil {
		return nil, err
	}
	return captData, nil
}
