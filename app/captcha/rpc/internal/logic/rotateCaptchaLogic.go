package logic

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	images "github.com/wenlng/go-captcha-assets/resources/images_v2"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/rotate"
	"github.com/zeromicro/go-zero/core/logx"
	"graph-med/app/captcha/rpc/internal/svc"
	"graph-med/app/captcha/rpc/pd"
	"log"
)

type RotateCaptchaLogic struct {
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	rotateCapt rotate.Captcha
	logx.Logger
}

func NewRotateCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RotateCaptchaLogic {

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

	rotateCapt := builder.Make()

	return &RotateCaptchaLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		rotateCapt: rotateCapt,
		Logger:     logx.WithContext(ctx),
	}
}

func (l *RotateCaptchaLogic) RotateCaptcha(in *pd.RotateCaptchaReq) (*pd.RotateCaptchaResp, error) {
	captData, err := l.rotateCapt.Generate()
	if err != nil {
		return nil, err
	}

	captchaKey := uuid.New().String()
	angle := captData.GetData().Angle

	graphMedCaptchaKey := fmt.Sprintf("%s%v", cacheGraphMedCaptchaRotatePrefix, captchaKey)
	angleVlaue := fmt.Sprintf("%v", angle)
	err = l.svcCtx.RedisClient.Setex(graphMedCaptchaKey, angleVlaue, captchaExpire)
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

	return &pd.RotateCaptchaResp{
		CaptchaId:   captchaKey,
		ImageBase64: imageBase64,
		ThumbBase64: thumbBase64,
		ParentSize:  int32(captData.GetData().ParentWidth),
		ChildSize:   int32(captData.GetData().Width),
	}, nil
}
