package logic

import (
	"context"
	"encoding/json"
	"fmt"
	mcpClient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pkg/errors"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	chatModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"graph-med/app/chat/model"
	"graph-med/pkg/xerr"
	"io"
	"os"
	"sync"
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
	var mcpIds []string
	chatType, err := l.svcCtx.ChatTypeModel.FindOneByTypeId(l.ctx, session.TypeId)
	if err != nil || chatType == nil {
		modelName = DefaultModelName
		mcpIds = []string{}
	} else {
		modelName = chatType.ModelName
		mcpIds = chatType.McpIds
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

	// 加载MCP服务
	mcpServices, err := l.svcCtx.McpServiceModel.FindAllByMcpIds(l.ctx, mcpIds)
	if err != nil {
		return err
	}

	// 初始化MCP client
	var wg sync.WaitGroup
	clientList := make([]mcpClient.MCPClient, 0)
	clientToolsList := make([][]*chatModel.Tool, 0)
	for _, mcpService := range mcpServices {
		wg.Add(1)
		go func(mcpService *model.McpService) {
			defer wg.Done()
			// 初始化MCP
			client, clientTools, err := InitMcpClient(l.ctx, mcpService)
			if err != nil {
				return
			}
			clientToolsList = append(clientToolsList, clientTools)
			clientList = append(clientList, client)
		}(mcpService)
	}
	wg.Wait()

	// 初始化funcNameToMcpClient
	funcNameToMcpClient := make(map[string]*mcpClient.MCPClient)
	tools := make([]*chatModel.Tool, 0)
	for i, client := range clientList {
		tools = append(tools, clientToolsList[i]...)
		for _, tool := range clientToolsList[i] {
			funcNameToMcpClient[tool.Function.Name] = &client
		}
	}

	// 开启对话
	resultChan := make(chan model.MessageTurn, 10)
	stopChan := make(chan struct{}, 1)
	go StartChatStream(l.ctx, l.client, modelName, histories, tools, funcNameToMcpClient, resultChan, stopChan)

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

func StartChatStream(
	ctx context.Context,
	client *arkruntime.Client,
	modelName string,
	histories []*chatModel.ChatCompletionMessage,
	tools []*chatModel.Tool,
	funcNameToMcpClient map[string]*mcpClient.MCPClient,
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

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		logx.WithContext(ctx).Errorf("stream api error: %v", err)
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
			logx.WithContext(ctx).Errorf("stream api error: %v", err)
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
						result, err := FunctionCall(ctx, funcNameToMcpClient, functionName, functionArguments)

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

					StartChatStream(ctx, client, modelName, histories, tools, funcNameToMcpClient, resultChan, stopChan)
					return
				}
				break
			}
		}
	}

	stopChan <- struct{}{}
	return
}

func FunctionCall(ctx context.Context, funcNameToMcpClient map[string]*mcpClient.MCPClient, functionName, functionArguments string) (string, error) {
	client, ok := funcNameToMcpClient[functionName]
	if !ok {
		return "函数不存在", errors.New("函数不存在")
	}

	mapArguments := make(map[string]any)
	err := json.Unmarshal([]byte(functionArguments), &mapArguments)
	if err != nil {
		return "函数调用失败", err
	}

	request := mcp.CallToolRequest{}
	request.Params.Name = functionName
	request.Params.Arguments = mapArguments
	result, err := (*client).CallTool(ctx, request)
	if err != nil || result.IsError || len(result.Content) == 0 {
		return "函数调用失败", errors.New("函数调用失败")
	}

	return result.Content[0].(mcp.TextContent).Text, nil
}

func InitMcpClient(ctx context.Context, service *model.McpService) (client mcpClient.MCPClient, tools []*chatModel.Tool, err error) {
	// Initialize
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "mcp-client",
		Version: "1.0.0",
	}

	// 根据 service.Type 初始化 client
	if service.Type == "sse" {
		sseClient, err := mcpClient.NewSSEMCPClient(service.BaseUrl)
		if err != nil {
			return nil, nil, err
		}

		if err := sseClient.Start(ctx); err != nil {
			return nil, nil, err
		}

		client = sseClient
	} else if service.Type == "stdio" {
		client, err = mcpClient.NewStdioMCPClient(service.Command, []string{})
		if err != nil {
			return nil, nil, err
		}
	}

	result, err := client.Initialize(ctx, initRequest)
	if err != nil {
		return nil, nil, err
	}

	if result.ServerInfo.Name != service.Name {
		return nil, nil, errors.New("mcp server name not match")
	}

	tools = make([]*chatModel.Tool, 0)
	request := mcp.ListToolsRequest{}
	toolListResult, err := client.ListTools(ctx, request)
	if err != nil {
		fmt.Println("list tools error", err)
		return
	}

	for _, tool := range toolListResult.Tools {
		tools = append(tools, &chatModel.Tool{
			Type: "function",
			Function: &chatModel.FunctionDefinition{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.RawInputSchema,
			},
		})
	}

	return client, tools, nil
}
