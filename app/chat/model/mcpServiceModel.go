package genModel

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ McpServiceModel = (*customMcpServiceModel)(nil)

type (
	// McpServiceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMcpServiceModel.
	McpServiceModel interface {
		mcpServiceModel
	}

	customMcpServiceModel struct {
		*defaultMcpServiceModel
	}
)

// NewMcpServiceModel returns a model for the database table.
func NewMcpServiceModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) McpServiceModel {
	return &customMcpServiceModel{
		defaultMcpServiceModel: newMcpServiceModel(conn, c, opts...),
	}
}
