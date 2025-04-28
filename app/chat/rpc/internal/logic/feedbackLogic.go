package logic

import (
	"context"
	"graph-med/app/chat/model"
	"graph-med/pkg/xerr"

	"graph-med/app/chat/rpc/internal/svc"
	"graph-med/app/chat/rpc/pd"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrChatMessageNoExists = xerr.NewErrMsg("消息不存在")
	ErrChatSessionNotAuth  = xerr.NewErrMsg("没有权限")
)

type FeedbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedbackLogic {
	return &FeedbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 对话反馈
func (l *FeedbackLogic) Feedback(in *pd.FeedbackReq) (*pd.FeedbackResp, error) {

	// 检查session是否存在
	chatSession, err := l.svcCtx.ChatSessionModel.FindOneBySessionId(l.ctx, in.SessionId)
	if err != nil {
		return nil, err
	}
	if chatSession == nil {
		return nil, ErrChatSessionNoExists
	}

	if chatSession.UserId != in.UserId {
		return nil, ErrChatSessionNotAuth
	}

	// 检查message是否存在
	chatMessage, err := l.svcCtx.ChatMessageModel.FindOneBySessionIdAndMsgId(l.ctx, in.SessionId, in.MessageId)
	if err != nil {
		return nil, err
	}
	if chatMessage == nil {
		return nil, ErrChatMessageNoExists
	}

	feedback := &model.ChatFeedback{
		UserId:    in.UserId,
		SessionId: in.SessionId,
		MessageId: in.MessageId,
		Feedback:  in.Feedback,
	}

	err = l.svcCtx.ChatFeedbackModel.Insert(l.ctx, feedback)
	if err != nil {
		return nil, err
	}
	return &pd.FeedbackResp{}, nil
}
