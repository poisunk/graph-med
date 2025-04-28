package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *customChatTypeModel) FindOneByTypeId(ctx context.Context, typeId string) (*ChatType, error) {
	var data *ChatType

	chatTypeCacheKey := fmt.Sprintf("%s%v", prefixChatTypeCacheKey, typeId)
	err := m.conn.FindOne(ctx, chatTypeCacheKey, data, bson.M{"typeId": typeId})
	switch err {
	case nil:
		return data, nil
	case mon.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
