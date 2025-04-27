package genModel

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChatTypeModel = (*customChatTypeModel)(nil)

type (
	// ChatTypeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatTypeModel.
	ChatTypeModel interface {
		chatTypeModel
	}

	customChatTypeModel struct {
		*defaultChatTypeModel
	}
)

// NewChatTypeModel returns a model for the database table.
func NewChatTypeModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ChatTypeModel {
	return &customChatTypeModel{
		defaultChatTypeModel: newChatTypeModel(conn, c, opts...),
	}
}
