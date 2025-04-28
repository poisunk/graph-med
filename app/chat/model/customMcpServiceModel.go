package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *customMcpServiceModel) FindAllByMcpIds(ctx context.Context, mcpIds []string) ([]*McpService, error) {
	var data []*McpService

	err := m.conn.Find(ctx, &data, bson.M{"mcpIds": bson.M{"$in": mcpIds}})
	switch err {
	case nil:
		return data, nil
	case mon.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
