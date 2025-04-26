package user

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"graph-med/app/usercenter/model"
	"graph-med/app/usercenter/rpc/usercenter"

	"graph-med/app/usercenter/api/internal/svc"
	"graph-med/app/usercenter/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// register
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	registerResp, err := l.svcCtx.UsercenterRpc.Register(l.ctx, &usercenter.RegisterReq{
		Email:    req.Email,
		Password: req.Password,
		AuthKey:  req.Email,
		AuthType: model.UserAuthTypeSystem,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}

	resp = &types.RegisterResp{}
	_ = copier.Copy(resp, registerResp)

	return resp, nil
}
