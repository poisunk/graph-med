package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *customChatMessageModel) FindAllBySessionId(ctx context.Context, sessionId string) ([]*ChatMessage, error) {
	var data []*ChatMessage

	err := m.conn.Find(ctx, &data, bson.M{"sessionId": sessionId})
	switch err {
	case nil:
		return data, nil
	case mon.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customChatMessageModel) InsertAll(ctx context.Context, messages []*ChatMessage) error {
	var data []any

	for _, message := range messages {
		data = append(data, message)
	}

	_, err := m.conn.InsertMany(ctx, data)

	return err
}

func (m *customChatMessageModel) FindOneBySessionIdAndMsgId(ctx context.Context, sessionId string, msgId int64) (*ChatMessage, error) {
	var data ChatMessage

	err := m.conn.FindOne(ctx, &data, bson.M{"sessionId": sessionId, "messageId": msgId})
	switch err {
	case nil:
		return &data, nil
	case mon.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
