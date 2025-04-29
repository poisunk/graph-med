package chat

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"graph-med/app/chat/rpc/chat"
	"graph-med/pkg/ctxdata"

	"graph-med/app/chat/api/internal/svc"
	"graph-med/app/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatCompletionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发起对话
func NewChatCompletionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatCompletionLogic {
	return &ChatCompletionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatCompletionLogic) ChatCompletion(req *types.ChatCompletionReq) (chan string, error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)

	client, err := l.svcCtx.ChatRpc.ChatCompletion(l.ctx, &chat.ChatCompletionReq{
		UserId:          userId,
		SessionId:       req.SessionId,
		Prompt:          req.Prompt,
		ParentMessageId: req.ParentMessageId,
	})

	if err != nil {
		return nil, err
	}

	resultChan := make(chan string, 10)

	go func() {
		defer close(resultChan)
		for {
			recv, err := client.Recv()
			if err != nil {
				return
			}

			var resp types.ChatCompletionResp
			_ = copier.Copy(&resp, recv)

			respData, err := json.Marshal(resp)
			if err != nil {
				return
			}

			resultChan <- string(respData)
		}
	}()

	return resultChan, nil
}
