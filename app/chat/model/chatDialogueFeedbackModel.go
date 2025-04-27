package genModel

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChatDialogueFeedbackModel = (*customChatDialogueFeedbackModel)(nil)

type (
	// ChatDialogueFeedbackModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatDialogueFeedbackModel.
	ChatDialogueFeedbackModel interface {
		chatDialogueFeedbackModel
	}

	customChatDialogueFeedbackModel struct {
		*defaultChatDialogueFeedbackModel
	}
)

// NewChatDialogueFeedbackModel returns a model for the database table.
func NewChatDialogueFeedbackModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ChatDialogueFeedbackModel {
	return &customChatDialogueFeedbackModel{
		defaultChatDialogueFeedbackModel: newChatDialogueFeedbackModel(conn, c, opts...),
	}
}
