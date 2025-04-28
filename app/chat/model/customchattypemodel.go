package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *customChatTypeModel) FindOneByTypeId(ctx context.Context, typeId string) (*ChatType, error) {
	var data *ChatType

	err := m.conn.FindOne(ctx, data, bson.M{"typeId": typeId})
	switch err {
	case nil:
		return data, nil
	case mon.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
