package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"graph-med/internal/base/response"
	"graph-med/internal/service"
	"graph-med/internal/utils"
)

type ChatController struct {
	chatService *service.ChatService
}

func NewChatController(chatService *service.ChatService) *ChatController {
	return &ChatController{
		chatService: chatService,
	}
}

// Chat 对话接口
func (c *ChatController) Chat(ctx *gin.Context) {
	var req struct {
		ChatSessionID   string `json:"chat_session_id" binding:"required"`
		ParentMessageID int    `json:"parent_message_id"`
		Prompt          string `json:"prompt" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, err)
		return
	}

	// 检查会话是否存在
	exist, err := c.chatService.CheckChatSession(req.ChatSessionID)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	if !exist {
		response.Error(ctx, fmt.Errorf("会话不存在"))
		return
	}

	resultChan, err := c.chatService.UserQueryChat(
		req.ChatSessionID,
		req.ParentMessageID,
		req.Prompt,
	)

	if err != nil {
		response.Error(ctx, err)
		return
	}

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	for res := range resultChan {
		err = response.EventStream(ctx, res)
		if err != nil {
			fmt.Printf("write error: %v\n", err)
			return
		}
	}
}

// CreateSession 创建会话
func (c *ChatController) CreateSession(ctx *gin.Context) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		response.Unauthorized(ctx)
		return
	}

	session, err := c.chatService.CreateChatSession("default", userId.(string))
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, session)
}

// DeleteSession 删除会话
func (c *ChatController) DeleteSession(ctx *gin.Context) {
	sessionId := ctx.Param("id")
	if sessionId == "" {
		response.Error(ctx, fmt.Errorf("会话id不能为空"))
		return
	}

	valid := utils.IsValidUUID(sessionId)
	if !valid {
		response.Error(ctx, fmt.Errorf("会话id不合法"))
		return
	}

	err := c.chatService.DeleteChatSession(sessionId)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, nil)
}

// UpdateSession 更新会话
func (c *ChatController) UpdateSession(ctx *gin.Context) {
	sessionId := ctx.Param("id")
	if sessionId == "" {
		response.Error(ctx, fmt.Errorf("会话id不能为空"))
		return
	}

	var req struct {
		Title string `json:"title" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, err)
		return
	}

	valid := utils.IsValidUUID(sessionId)
	if !valid {
		response.Error(ctx, fmt.Errorf("会话id不合法"))
		return
	}

	err := c.chatService.UpdateChatSession(sessionId, req.Title)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, nil)
}

// Feedback 反馈
func (c *ChatController) Feedback(ctx *gin.Context) {
	var req struct {
		ChatSessionID string `json:"chat_session_id" binding:"required"`
		MessageID     int    `json:"message_id" binding:"required"`
		Feedback      string `json:"feedback" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, err)
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Unauthorized(ctx)
		return
	}

	err := c.chatService.Feedback(req.ChatSessionID, userID.(string), req.MessageID, req.Feedback)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}
