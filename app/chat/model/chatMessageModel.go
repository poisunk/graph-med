package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/mon"
)

var _ ChatMessageModel = (*customChatMessageModel)(nil)

type (
	// ChatMessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatMessageModel.
	ChatMessageModel interface {
		chatMessageModel
		FindAllBySessionId(ctx context.Context, sessionId string) ([]*ChatMessage, error)
		InsertAll(ctx context.Context, messages []*ChatMessage) error
		FindOneBySessionIdAndMsgId(ctx context.Context, sessionId string, msgId int64) (*ChatMessage, error)
		DeleteMessageInSessionWithIds(ctx context.Context, sessionId string, msgIds []int64) error
	}

	customChatMessageModel struct {
		*defaultChatMessageModel
	}
)

// NewChatMessageModel returns a model for the mongo.
func NewChatMessageModel(url, db, collection string) ChatMessageModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customChatMessageModel{
		defaultChatMessageModel: newDefaultChatMessageModel(conn),
	}
}
