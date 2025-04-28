package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/monc"
)

var _ ChatTypeModel = (*customChatTypeModel)(nil)

type (
	// ChatTypeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatTypeModel.
	ChatTypeModel interface {
		chatTypeModel
		FindOneByTypeId(ctx context.Context, typeId string) (*ChatType, error)
	}

	customChatTypeModel struct {
		*defaultChatTypeModel
	}
)

// NewChatTypeModel returns a model for the mongo.
func NewChatTypeModel(url, db, collection string, c cache.CacheConf) ChatTypeModel {
	conn := monc.MustNewModel(url, db, collection, c)
	return &customChatTypeModel{
		defaultChatTypeModel: newDefaultChatTypeModel(conn),
	}
}
