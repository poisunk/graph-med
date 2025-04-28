package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	chatModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"graph-med/app/chat/model"
	"os"
	"time"

	"graph-med/app/chat/rpc/internal/svc"
	"graph-med/app/chat/rpc/pd"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegenerateChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	client *arkruntime.Client

	logx.Logger
}

func NewRegenerateChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegenerateChatLogic {
	client := arkruntime.NewClientWithApiKey(
		os.Getenv("LLM_API_KEY"),
		arkruntime.WithTimeout(30*time.Minute),
	)

	return &RegenerateChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		client: client,
		Logger: logx.WithContext(ctx),
	}
}

// 重新发起对话
func (l *RegenerateChatLogic) RegenerateChat(in *pd.ChatCompletionReq, stream pd.Chat_RegenerateChatServer) error {
	// 判断session是否存在
	session, err := l.svcCtx.ChatSessionModel.FindOneBySessionId(l.ctx, in.SessionId)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.Wrapf(ErrChatSessionNoExists, "session_id:%s", in.SessionId)
	}

	// 存在就查找ChatType
	var modelName string
	chatType, err := l.svcCtx.ChatTypeModel.FindOneByTypeId(l.ctx, session.TypeId)
	if !errors.Is(err, model.ErrNotFound) && err != nil {
		return err
	}
	if errors.Is(err, model.ErrNotFound) || chatType == nil {
		modelName = DefaultModelName
	} else {
		modelName = chatType.ModelName
	}

	// 组装userMessage
	userMessageId := in.ParentMessageId + 1
	userMessage := &model.ChatMessage{
		SessionId:       in.SessionId,
		MessageId:       userMessageId,
		ParentMessageId: in.ParentMessageId,
		Role:            chatModel.ChatMessageRoleUser,
		Content:         in.Prompt,
	}

	// 查找历史对话
	chatHistories, err := l.svcCtx.ChatMessageModel.FindAllBySessionId(l.ctx, in.SessionId)
	if err != nil {
		return err
	}

	deleteIds := make([]int64, 0)
	var histories []*chatModel.ChatCompletionMessage
	for _, chatHistory := range chatHistories {
		if chatHistory.MessageId > in.ParentMessageId {
			deleteIds = append(deleteIds, chatHistory.MessageId)
			continue
		}
		histories = append(histories, &chatModel.ChatCompletionMessage{
			Role: chatHistory.Role,
			Content: &chatModel.ChatCompletionMessageContent{
				StringValue: volcengine.String(chatHistory.Content),
			},
		})
	}
	histories = append(histories, &chatModel.ChatCompletionMessage{
		Role: "user",
		Content: &chatModel.ChatCompletionMessageContent{
			StringValue: volcengine.String(in.Prompt),
		},
	})

	// 删除多余message
	err = l.svcCtx.ChatMessageModel.DeleteMessageInSessionWithIds(l.ctx, in.SessionId, deleteIds)
	if err != nil {
		return err
	}

	// TODO: MCP

	// 开启对话
	resultChan := make(chan model.MessageTurn, 10)
	stopChan := make(chan struct{}, 1)
	go StartChatStream(l.ctx, l.client, modelName, histories, nil, resultChan, stopChan)

	created := time.Now().Unix()
	messageTurn := &model.MessageTurn{}
	messageTurnList := make([]*model.MessageTurn, 0)
	for {
		select {
		case turn := <-resultChan:
			resp := pd.ChatCompletionResp{
				Choices: []*pd.ChatCompletionResp_Choice{
					{
						Delta: &pd.ChatCompletionResp_Choice_Delta{
							Role:    chatModel.ChatMessageRoleAssistant,
							Content: turn.Content,
						},
					},
				},
				PromptTokenUsage: 1,
				Created:          created,
			}

			// 如果消息的类型切换
			if messageTurn.Type != turn.Type {
				if messageTurn.Content != "" {
					messageTurnList = append(messageTurnList, messageTurn)
				}
				messageTurn = &model.MessageTurn{
					Type: turn.Type,
				}
			}

			// 如果是function call类型
			if turn.Type == MessageContentFunctionCallType {
				resp.Choices[0].ToolCalls = []*pd.ChatCompletionResp_Choice_ToolCall{
					{
						Name:      turn.FunctionName,
						Arguments: turn.FunctionArgs,
					},
				}
				messageTurn.FunctionName = turn.FunctionName
				messageTurn.FunctionArgs = turn.FunctionArgs
			}

			messageTurn.Content += turn.Content

			err := stream.Send(&resp)
			if err != nil {
				return err
			}
		case <-stopChan:
			close(resultChan)
			go func() {
				messageTurnList = append(messageTurnList, messageTurn)

				ctx := context.Background()
				assistantMessageId := userMessageId + 1
				assistantMessage := &model.ChatMessage{
					SessionId:       in.SessionId,
					MessageId:       assistantMessageId,
					ParentMessageId: userMessageId,
					Role:            chatModel.ChatMessageRoleAssistant,
					Turns:           messageTurnList,
				}
				err := l.svcCtx.ChatMessageModel.InsertAll(ctx, []*model.ChatMessage{userMessage, assistantMessage})
				if err != nil {
					l.Logger.WithContext(ctx).Errorf("insert chat ialogue error: %v", err)
				}
			}()

			return nil
		}
	}
}
