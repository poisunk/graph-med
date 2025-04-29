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

type FeedbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 对话反馈
func NewFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedbackLogic {
	return &FeedbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FeedbackLogic) Feedback(req *types.FeedbackReq) (resp *types.FeedbackResp, err error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)

	result, err := l.svcCtx.ChatRpc.Feedback(l.ctx, &chat.FeedbackReq{
		UserId:    userId,
		SessionId: req.SsessionId,
		MessageId: req.MessageId,
		Feedback:  req.Feedback,
	})

	if err != nil {
		return nil, err
	}

	resp = &types.FeedbackResp{}
	_ = copier.Copy(resp, result)

	return
}
