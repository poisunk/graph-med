package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *customChatSessionModel) FindOneBySessionId(ctx context.Context, sessionId string) (*ChatSession, error) {
	var data ChatSession

	err := m.conn.FindOne(ctx, &data, bson.M{"sessionId": sessionId})
	switch err {
	case nil:
		return &data, nil
	case mon.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
