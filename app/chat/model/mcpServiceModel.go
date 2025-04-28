package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/monc"
)

var _ McpServiceModel = (*customMcpServiceModel)(nil)

type (
	// McpServiceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMcpServiceModel.
	McpServiceModel interface {
		mcpServiceModel
		FindAllByMcpIds(ctx context.Context, mcpIds []string) ([]*McpService, error)
	}

	customMcpServiceModel struct {
		*defaultMcpServiceModel
	}
)

// NewMcpServiceModel returns a model for the mongo.
func NewMcpServiceModel(url, db, collection string, c cache.CacheConf) McpServiceModel {
	conn := monc.MustNewModel(url, db, collection, c)
	return &customMcpServiceModel{
		defaultMcpServiceModel: newDefaultMcpServiceModel(conn),
	}
}
