package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	chatModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"graph-med/app/chat/model"
	"graph-med/pkg/xerr"
	"io"
	"os"
	"time"

	"graph-med/app/chat/rpc/internal/svc"
	"graph-med/app/chat/rpc/pd"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrChatSessionNoExists = xerr.NewErrMsg("聊天会话不存在")

	MessageContentThinkingType     = "thinking"
	MessageContentTextType         = "text"
	MessageContentFunctionCallType = "function_call"

	DefaultModelName = "doubao-1-5-lite-32k-250115"
)

type ChatCompletionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	client *arkruntime.Client

	logx.Logger
}

func NewChatCompletionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatCompletionLogic {
	client := arkruntime.NewClientWithApiKey(
		os.Getenv("LLM_API_KEY"),
		arkruntime.WithTimeout(30*time.Minute),
	)

	return &ChatCompletionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		client: client,
		Logger: logx.WithContext(ctx),
	}
}

// ChatCompletion 发起对话
func (l *ChatCompletionLogic) ChatCompletion(in *pd.ChatCompletionReq, stream pd.Chat_ChatCompletionServer) error {

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

	var histories []*chatModel.ChatCompletionMessage
	for _, chatHistory := range chatHistories {
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

	// TODO: MCP

	// 开启对话
	resultChan := make(chan model.MessageTurn, 10)
	stopChan := make(chan struct{}, 1)
	go l.startChatStream(modelName, histories, nil, resultChan, stopChan)

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

func (l *ChatCompletionLogic) startChatStream(
	modelName string,
	histories []*chatModel.ChatCompletionMessage,
	tools []*chatModel.Tool,
	resultChan chan<- model.MessageTurn,
	stopChan chan<- struct{},
) {
	content := ""

	resp := model.MessageTurn{
		Content: "",
	}

	req := chatModel.ChatCompletionRequest{
		Model:    modelName,
		Messages: histories,
		Tools:    tools,
	}

	stream, err := l.client.CreateChatCompletionStream(l.ctx, req)
	if err != nil {
		l.Logger.WithContext(l.ctx).Errorf("stream api error: %v", err)
		stopChan <- struct{}{}
		return
	}
	defer stream.Close()

	finalToolCalls := make(map[int]*chatModel.ToolCall)
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			l.Logger.WithContext(l.ctx).Errorf("stream api error: %v", err)
			stopChan <- struct{}{}
			return
		}

		if len(recv.Choices) > 0 {

			for _, tool_call := range recv.Choices[0].Delta.ToolCalls {
				if _, ok := finalToolCalls[*tool_call.Index]; !ok {
					finalToolCalls[*tool_call.Index] = tool_call
				}
				finalToolCalls[*tool_call.Index].Function.Arguments += tool_call.Function.Arguments
			}

			if recv.Choices[0].Delta.ReasoningContent != nil {
				resp.Content = *recv.Choices[0].Delta.ReasoningContent
				resp.Type = MessageContentThinkingType
				resultChan <- resp
				content = content + *recv.Choices[0].Delta.ReasoningContent
			}

			if recv.Choices[0].Delta.Content != "" {
				resp.Content = recv.Choices[0].Delta.Content
				resp.Type = MessageContentTextType
				resultChan <- resp
				content = content + recv.Choices[0].Delta.Content
			}

			if recv.Choices[0].FinishReason != "" {
				if recv.Choices[0].FinishReason == chatModel.FinishReasonToolCalls {
					// 如果是工具调用
					toolCalls := make([]*chatModel.ToolCall, len(finalToolCalls))
					for i, toolCall := range finalToolCalls {
						toolCalls[i] = toolCall
					}
					histories = append(histories, &chatModel.ChatCompletionMessage{
						Role: chatModel.ChatMessageRoleAssistant,
						Content: &chatModel.ChatCompletionMessageContent{
							StringValue: volcengine.String(content),
						},
						ToolCalls: toolCalls,
					})

					// 如果有 处理工具调用
					for _, toolCall := range finalToolCalls {
						functionName := toolCall.Function.Name
						functionArguments := toolCall.Function.Arguments
						toolCallID := toolCall.ID
						//l.Logger.WithContext(ctx).Infof("调用工具:%s, %s, %s", functionName, functionArguments, toolCallID)
						// 调用工具
						result, err := l.functionCall(functionName, functionArguments)

						functionCallResult := result
						if err != nil {
							functionCallResult = ""
						}
						resp.Content = functionCallResult
						resp.Type = MessageContentFunctionCallType
						resp.FunctionName = functionName
						resp.FunctionArgs = functionArguments
						resultChan <- resp
						if err != nil {
							return
						}

						histories = append(histories, &chatModel.ChatCompletionMessage{
							Role: chatModel.ChatMessageRoleTool,
							Content: &chatModel.ChatCompletionMessageContent{
								StringValue: volcengine.String(result),
							},
							Name:       &functionName,
							ToolCallID: toolCallID,
						})
					}

					l.startChatStream(modelName, histories, tools, resultChan, stopChan)
					return
				}
				break
			}
		}
	}

	stopChan <- struct{}{}
	return
}

func (l *ChatCompletionLogic) functionCall(functionName, functionArguments string) (string, error) {
	// TODO
	return "", nil
}
