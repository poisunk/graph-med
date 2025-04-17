package repository

import (
	"errors"
	"graph-med/internal/base/data"
	"graph-med/internal/model"
)

type ChatRepository struct {
	data *data.Data
}

func NewChatRepository(data *data.Data) *ChatRepository {
	return &ChatRepository{data: data}
}

// CreateChatMessage 创建对话消息
func (r *ChatRepository) CreateChatMessage(Message *model.ChatMessage) error {
	_, err := r.data.DB.Insert(Message)
	return err
}

// CreateChatMessages 批量创建聊天会话
func (r *ChatRepository) CreateChatMessages(messages ...*model.ChatMessage) error {
	_, err := r.data.DB.Insert(messages)
	return err
}

// DeleteChatMessage 根据ID删除聊天会话
func (r *ChatRepository) DeleteChatMessage(id int) error {
	_, err := r.data.DB.ID(id).Delete(&model.ChatMessage{})
	return err
}

// CountMessageBySessionID 统计会话中的消息数量
func (r *ChatRepository) CountMessageBySessionID(sessionID string) (int64, error) {
	return r.data.DB.Where("chat_session_id = ?", sessionID).Count(&model.ChatMessage{})
}

// GetChatMessageBySessionID 获取会话消息
func (r *ChatRepository) GetChatMessageBySessionID(sessionID string, limit int) ([]*model.ChatMessage, error) {
	messages := make([]*model.ChatMessage, 0)
	err := r.data.DB.Where("chat_session_id = ?", sessionID).Limit(limit).Desc("created_at").Find(&messages)
	return messages, err
}

// CreateChatSession 创建聊天会话
func (r *ChatRepository) CreateChatSession(session *model.ChatSession) error {
	_, err := r.data.DB.Insert(session)
	return err
}

// FindChatSession 根据会话ID查找会话
func (r *ChatRepository) FindChatSession(sessionID string) (*model.ChatSession, error) {
	session := &model.ChatSession{}
	exists, err := r.data.DB.Where("session_id = ?", sessionID).Get(session)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("session not found")
	}
	return session, nil
}

// DeleteChatSession 根据会话ID删除会话
func (r *ChatRepository) DeleteChatSession(sessionID string) error {
	_, err := r.data.DB.Where("session_id = ?", sessionID).Delete(&model.ChatSession{})
	return err
}

// DeleteChatMessageBySessionID 根据会话ID删除会话消息
func (r *ChatRepository) DeleteChatMessageBySessionID(sessionID string) error {
	_, err := r.data.DB.Where("chat_session_id = ?", sessionID).Delete(&model.ChatMessage{})
	return err
}

// UpdateChatSession 更新会话
func (r *ChatRepository) UpdateChatSession(session *model.ChatSession) error {
	_, err := r.data.DB.ID(session.ID).Update(session)
	return err
}

// CreateChatFeedback 创建反馈
func (r *ChatRepository) CreateChatFeedback(feedback *model.ChatFeedback) error {
	_, err := r.data.DB.Insert(feedback)
	return err
}

// GetChatMessageBySessionID 获取对话类型
func (r *ChatRepository) GetChatTypeBySessionID(sessionID string) (*model.ChatType, error) {
	session, err := r.FindChatSession(sessionID)
	if err != nil {
		return nil, err
	}

	chatType := &model.ChatType{}
	exists, err := r.data.DB.Where("type_id = ?", session.TypeID).Get(chatType)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("session not found")
	}
	return chatType, nil
}
