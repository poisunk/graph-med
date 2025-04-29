package chat

import (
	"context"
	"github.com/jinzhu/copier"
	"graph-med/app/chat/rpc/chat"
	"graph-med/pkg/ctxdata"

	"graph-med/app/chat/api/internal/svc"
	"graph-med/app/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateChatSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建对话session
func NewCreateChatSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateChatSessionLogic {
	return &CreateChatSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateChatSessionLogic) CreateChatSession(req *types.CreateChatSessionReq) (resp *types.CreateChatSessionResp, err error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)

	result, err := l.svcCtx.ChatRpc.CreateChatSession(l.ctx, &chat.CreateChatSessionReq{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}

	resp = &types.CreateChatSessionResp{}
	_ = copier.Copy(resp, result)

	return
}
