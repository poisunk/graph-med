package logic

import (
	"context"
	"github.com/google/uuid"
	"graph-med/app/chat/model"
	"graph-med/app/chat/rpc/internal/svc"
	"graph-med/app/chat/rpc/pd"
	"graph-med/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrChatSessionInsertError = xerr.NewErrMsg("创建聊天会话失败")
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
	sessionId := uuid.New().String()

	err := l.svcCtx.ChatSessionModel.Insert(l.ctx, &model.ChatSession{
		SessionId: sessionId,
		TypeId:    "default",
		UserId:    in.UserId,
	})

	if err != nil {
		return nil, err
	}

	chatSession, err := l.svcCtx.ChatSessionModel.FindOneBySessionId(l.ctx, sessionId)

	if err != nil {
		return nil, err
	}

	if chatSession == nil {
		return nil, ErrChatSessionInsertError
	}

	return &pd.CreateChatSessionResp{
		SessionId: sessionId,
		CreatedAt: chatSession.CreateAt.Format("2006-01-02 15:04:05"),
	}, nil
}
