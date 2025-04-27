package genModel

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChatDialogueModel = (*customChatDialogueModel)(nil)

type (
	// ChatDialogueModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatDialogueModel.
	ChatDialogueModel interface {
		chatDialogueModel
	}

	customChatDialogueModel struct {
		*defaultChatDialogueModel
	}
)

// NewChatDialogueModel returns a model for the database table.
func NewChatDialogueModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ChatDialogueModel {
	return &customChatDialogueModel{
		defaultChatDialogueModel: newChatDialogueModel(conn, c, opts...),
	}
}
