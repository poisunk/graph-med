package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"graph-med/internal/base/constant"
	"graph-med/internal/model"
	"graph-med/internal/repository"
	"graph-med/internal/schema"
	"graph-med/internal/utils"
	"io"
	"os"
	"sync"
	"time"

	mcpClient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	chatModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type ChatService struct {
	client   *arkruntime.Client
	chatRepo *repository.ChatRepository
	mcpRepo  *repository.McpRepository
	// mcpID to client
	mcpServiceCache map[string]mcpClient.MCPClient
	// functionName to mcpID
	functionStore map[string]string
}

func NewChatService(chatRepo *repository.ChatRepository, mcpRepo *repository.McpRepository) *ChatService {
	client := arkruntime.NewClientWithApiKey(
		os.Getenv("LLM_API_KEY"),
		arkruntime.WithTimeout(30*time.Minute),
	)

	return &ChatService{
		client:          client,
		chatRepo:        chatRepo,
		mcpRepo:         mcpRepo,
		mcpServiceCache: make(map[string]mcpClient.MCPClient),
		functionStore:   make(map[string]string),
	}
}

// Chat 对话
func (c *ChatService) ChatStream(
	chatSessionId string,
	chatType *model.ChatType,
	userMsg *model.ChatMessage,
	histories []*chatModel.ChatCompletionMessage,
) (<-chan string, error) {
	if !utils.IsValidUUID(chatSessionId) {
		return nil, fmt.Errorf("invalid api session id")
	}

	// 创建会话
	resultChan := make(chan string)

	go func() {
		defer func() {
			resultChan <- "[DONE]"
			close(resultChan)
		}()
		ctx := context.Background()
		defer ctx.Done()

		// 获取会话对应的mcp services
		mcpServices, err := c.mcpRepo.GetMcpServerBySessionID(chatSessionId)
		if err != nil {
			return
		}

		saveMessages := []*model.ChatMessage{userMsg}

		assisParentMessageID := userMsg.MessageID
		assisMsgID := userMsg.MessageID + 1
		fmt.Println(assisMsgID, assisParentMessageID)
		createTime := time.Now().Unix()
		promptTokenUsage := 0
		for _, msg := range histories {
			promptTokenUsage += utils.NumTokens(*msg.Content.StringValue)
		}

		resp := &schema.ChatResp{
			Choices: []schema.ChatChoice{
				{
					Index: 0,
					Delta: schema.ChatDelta{
						Content: "",
						Type:    "text",
					},
				},
			},
			PromptTokenUsage: promptTokenUsage,
			Created:          createTime,
			MessageID:        assisMsgID,
			ParentID:         assisParentMessageID,
		}
		err = c.SendData(resultChan, resp)
		if err != nil {
			return
		}

		finished := false
		for !finished {
			finished = true
			// 发送消息
			content := ""
			resp.PromptTokenUsage = 0
			resp.ChunkTokenUsage = 1
			resp.MessageID = assisMsgID
			resp.ParentID = assisParentMessageID
			fmt.Println("histories", len(histories))
			req := chatModel.ChatCompletionRequest{
				Model:    chatType.ModelName,
				Messages: histories,
				Tools:    c.GetTools(ctx, mcpServices),
			}

			stream, err := c.client.CreateChatCompletionStream(ctx, req)
			if err != nil {
				fmt.Printf("stream api error: %v\n", err)
				return
			}
			defer stream.Close()

			final_tool_calls := make(map[int]*chatModel.ToolCall)
			for {
				recv, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Printf("Stream api error: %v", err)
					return
				}

				if len(recv.Choices) > 0 {

					for _, tool_call := range recv.Choices[0].Delta.ToolCalls {
						if _, ok := final_tool_calls[*tool_call.Index]; !ok {
							final_tool_calls[*tool_call.Index] = tool_call
						}
						final_tool_calls[*tool_call.Index].Function.Arguments += tool_call.Function.Arguments
					}

					if recv.Choices[0].Delta.ReasoningContent != nil {
						resp.Choices[0].Delta.Content = *recv.Choices[0].Delta.ReasoningContent
						resp.Choices[0].Delta.Type = "thinking"
						err = c.SendData(resultChan, resp)
						if err != nil {
							return
						}
						content = content + *recv.Choices[0].Delta.ReasoningContent
					}

					if recv.Choices[0].Delta.Content != "" {
						resp.Choices[0].Delta.Content = recv.Choices[0].Delta.Content
						resp.Choices[0].Delta.Type = "text"
						err = c.SendData(resultChan, resp)
						if err != nil {
							return
						}
						content = content + recv.Choices[0].Delta.Content
					}

					if recv.Choices[0].FinishReason != "" {
						resp.ChunkTokenUsage = 0
						resp.Choices[0].FinishReason = string(recv.Choices[0].FinishReason)

						fmt.Println("finish reason", recv.Choices[0].FinishReason)
						if recv.Choices[0].FinishReason == chatModel.FinishReasonToolCalls {
							// 如果是工具调用
							toolCalls := make([]*chatModel.ToolCall, len(final_tool_calls))
							for i, toolCall := range final_tool_calls {
								toolCalls[i] = toolCall
							}
							histories = append(histories, &chatModel.ChatCompletionMessage{
								Role: chatModel.ChatMessageRoleAssistant,
								Content: &chatModel.ChatCompletionMessageContent{
									StringValue: volcengine.String(content),
								},
								ToolCalls: toolCalls,
							})

							// 处理工具调用
							resp.Choices[0].ToolCalls = make([]schema.ToolCall, 1)
							for _, toolCall := range final_tool_calls {
								functionName := toolCall.Function.Name
								functionArguments := toolCall.Function.Arguments
								toolCallID := toolCall.ID
								fmt.Println("调用工具:", functionName, functionArguments, toolCallID)
								// 调用工具
								result, err := c.FunctionCall(ctx, functionName, functionArguments)

								// TODO
								functionCallResult := result
								if err != nil {
									functionCallResult = ""
								}
								resp.Choices[0].Delta.Content = functionCallResult
								resp.Choices[0].Delta.Type = "function_call"
								resp.Choices[0].ToolCalls[0].Name = functionName
								resp.Choices[0].ToolCalls[0].Arguments = functionArguments
								err = c.SendData(resultChan, resp)
								if err != nil {
									return
								}
								content += functionCallResult

								histories = append(histories, &chatModel.ChatCompletionMessage{
									Role: chatModel.ChatMessageRoleTool,
									Content: &chatModel.ChatCompletionMessageContent{
										StringValue: volcengine.String(result),
									},
									Name:       &functionName,
									ToolCallID: toolCallID,
								})
							}
							resp.Choices[0].ToolCalls = []schema.ToolCall{}
							finished = false
							break
						}

						saveMessages = append(saveMessages, &model.ChatMessage{
							MessageID:       assisMsgID,
							ParentMessageID: assisParentMessageID,
							ChatSessionID:   chatSessionId,
							Role:            "assistant",
							Content:         content,
							ModelName:       chatType.ModelName,
							Operator:        "system",
						})
						break
					}
				}
			}
		}

		// 保存消息
		err = c.chatRepo.CreateChatMessages(saveMessages...)
		if err != nil {
			return
		}
	}()

	return resultChan, nil
}

// UserQueryChat 用户咨询会话
func (c *ChatService) UserQueryChat(chatSessionId string, parentId int, prompt string) (<-chan string, error) {

	chatType, err := c.chatRepo.GetChatTypeBySessionID(chatSessionId)
	if err != nil {
		return nil, err
	}

	modelName := chatType.ModelName
	// 新建消息
	userMsgID, err := c.GetSessionNewMessageID(chatSessionId)
	if err != nil {
		return nil, err
	}

	userMsg := &model.ChatMessage{
		MessageID:       userMsgID,
		ParentMessageID: parentId,
		ChatSessionID:   chatSessionId,
		Role:            "user",
		Content:         prompt,
		ModelName:       modelName,
		Operator:        "system",
	}

	// 获取历史对话
	histories, err := c.GetSessionHistory(chatSessionId, parentId)
	if err != nil {
		return nil, err
	}
	histories = append(histories, &chatModel.ChatCompletionMessage{
		Role: chatModel.ChatMessageRoleUser,
		Content: &chatModel.ChatCompletionMessageContent{
			StringValue: volcengine.String(prompt),
		},
	})

	return c.ChatStream(chatSessionId, chatType, userMsg, histories)
}

// GetSessionNewMessageID 获取会话中的最新消息ID
func (c *ChatService) GetSessionNewMessageID(chatSessionID string) (int, error) {
	cnt, err := c.chatRepo.CountMessageBySessionID(chatSessionID)
	if err != nil {
		return 0, err
	}
	return int(cnt) + 1, nil
}

// GetSessionHistory 获取会话历史
func (c *ChatService) GetSessionHistory(chatSessionID string, parentId int) ([]*chatModel.ChatCompletionMessage, error) {
	msgList, err := c.chatRepo.GetChatMessageBySessionID(chatSessionID, constant.MAX_HISTORY_LENGTH)
	if err != nil {
		return nil, err
	}

	parentIdList := make([]int, 0, len(msgList))
	parentIdList = append(parentIdList, parentId)
	for _, msg := range msgList {
		if utils.InList(msg.MessageID, parentIdList) {
			parentIdList = append(parentIdList, msg.ParentMessageID)
		}
	}

	if len(parentIdList) == 0 {
		return []*chatModel.ChatCompletionMessage{}, nil
	}

	history := make([]*chatModel.ChatCompletionMessage, 0, len(msgList))
	for _, msg := range msgList {
		if !utils.InList(msg.MessageID, parentIdList) {
			continue
		}
		history = append(history, &chatModel.ChatCompletionMessage{
			Role: msg.Role,
			Content: &chatModel.ChatCompletionMessageContent{
				StringValue: volcengine.String(msg.Content),
			},
		})
	}
	return history, nil
}

// SendData 发送数据
func (c *ChatService) SendData(ch chan<- string, data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("json marshal error: %v\n", err)
		return err
	}
	ch <- string(jsonData)
	return nil
}

// CreateChatSession 创建消息
func (c *ChatService) CreateChatSession(typeId, userId string) (*schema.ChatSessionResp, error) {
	sessionId := utils.GenerateUUID()

	session := &model.ChatSession{
		TypeID:    typeId,
		SessionID: sessionId,
		UserID:    userId,
		Title:     "新会话",
	}
	err := c.chatRepo.CreateChatSession(session)
	if err != nil {
		return nil, err
	}

	ret := &schema.ChatSessionResp{
		ID:        sessionId,
		CreatedAt: session.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: session.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return ret, nil
}

// CheckChatSession 检查会话
func (c *ChatService) CheckChatSession(sessionID string) (bool, error) {
	_, err := c.chatRepo.FindChatSession(sessionID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteChatSession 删除会话
func (c *ChatService) DeleteChatSession(sessionID string) error {
	err := c.chatRepo.DeleteChatSession(sessionID)
	if err != nil {
		return err
	}

	err = c.chatRepo.DeleteChatMessageBySessionID(sessionID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateChatSession 更新会话
func (c *ChatService) UpdateChatSession(sessionID string, title string) error {
	session := &model.ChatSession{
		SessionID: sessionID,
		Title:     title,
	}
	err := c.chatRepo.UpdateChatSession(session)
	if err != nil {
		return err
	}
	return nil
}

// Feedback 反馈
func (c *ChatService) Feedback(chatSessionId string, userId string, messageId int, feedback string) error {
	err := c.chatRepo.CreateChatFeedback(&model.ChatFeedback{
		MessageID:     messageId,
		UserID:        userId,
		ChatSessionID: chatSessionId,
		Feedback:      feedback,
	})
	if err != nil {
		return err
	}
	return nil
}

// GetTools 获取工具
func (c *ChatService) GetTools(ctx context.Context, mcpServices []*model.McpService) []*chatModel.Tool {
	tools := make([]*chatModel.Tool, 0)

	var wg sync.WaitGroup

	wg.Add(len(mcpServices))

	for _, service := range mcpServices {
		go func(service *model.McpService) {
			defer wg.Done()

			client, err := c.initMcpClient(ctx, service)
			if err != nil {
				fmt.Println("init mcp client error", err)
				return
			}
			c.mcpServiceCache[service.McpID] = client

			request := mcp.ListToolsRequest{}
			result, err := client.ListTools(ctx, request)
			if err != nil {
				fmt.Println("list tools error", err)
				return
			}

			for _, tool := range result.Tools {
				tools = append(tools, &chatModel.Tool{
					Type: "function",
					Function: &chatModel.FunctionDefinition{
						Name:        tool.Name,
						Description: tool.Description,
						Parameters:  tool.RawInputSchema,
					},
				})

				c.functionStore[tool.Name] = service.McpID
			}
		}(service)
	}

	wg.Wait()

	return tools
}

// FunctionCall 函数调用
func (c *ChatService) FunctionCall(ctx context.Context, functionName, arguments string) (string, error) {
	mcpID, ok := c.functionStore[functionName]
	if !ok {
		return "函数不存在", errors.New("函数不存在")
	}
	client, ok := c.mcpServiceCache[mcpID]
	if !ok {
		return "函数不存在", errors.New("函数不存在")
	}

	mapArguments := make(map[string]any)
	err := json.Unmarshal([]byte(arguments), &mapArguments)
	if err != nil {
		return "函数调用失败", err
	}

	request := mcp.CallToolRequest{}
	request.Params.Name = functionName
	request.Params.Arguments = mapArguments
	result, err := client.CallTool(ctx, request)
	if err != nil || result.IsError || len(result.Content) == 0 {
		return "函数调用失败", errors.New("函数调用失败")
	}

	return result.Content[0].(mcp.TextContent).Text, nil
}

// initMcpClient
func (c *ChatService) initMcpClient(ctx context.Context, service *model.McpService) (client mcpClient.MCPClient, err error) {
	var mcpArgs map[string]string
	err = json.Unmarshal([]byte(service.Args), &mcpArgs)
	if err != nil {
		return nil, err
	}

	// Initialize
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "mcp-client",
		Version: "1.0.0",
	}

	// 根据 service.Type 初始化 client
	if service.Type == "sse" {
		sseClient, err := mcpClient.NewSSEMCPClient(mcpArgs["baseURL"] + "/sse")
		if err != nil {
			return nil, err
		}

		if err := sseClient.Start(ctx); err != nil {
			return nil, err
		}

		client = sseClient
	} else if service.Type == "stdio" {
		client, err = mcpClient.NewStdioMCPClient(mcpArgs["command"], []string{})
		if err != nil {
			return nil, err
		}
	}

	result, err := client.Initialize(ctx, initRequest)
	if err != nil {
		return nil, err
	}

	if result.ServerInfo.Name != service.Name {
		return nil, errors.New("mcp server name not match")
	}

	return client, nil
}
