package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/mon"
)

var _ ChatSessionModel = (*customChatSessionModel)(nil)

type (
	// ChatSessionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatSessionModel.
	ChatSessionModel interface {
		chatSessionModel
		FindOneBySessionId(ctx context.Context, sessionId string) (*ChatSession, error)
	}

	customChatSessionModel struct {
		*defaultChatSessionModel
	}
)

// NewChatSessionModel returns a model for the mongo.
func NewChatSessionModel(url, db, collection string) ChatSessionModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customChatSessionModel{
		defaultChatSessionModel: newDefaultChatSessionModel(conn),
	}
}
