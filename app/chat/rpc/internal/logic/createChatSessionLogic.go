package logic

import (
	"context"

	"graph-med/app/chat/rpc/internal/svc"
	"graph-med/app/chat/rpc/pd"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateChatSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateChatSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateChatSessionLogic {
	return &CreateChatSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateChatSessionLogic) CreateChatSession(in *pd.CreateChatSessionReq) (*pd.CreateChatSessionResp, error) {
	// todo: add your logic here and delete this line

	return &pd.CreateChatSessionResp{}, nil
}
