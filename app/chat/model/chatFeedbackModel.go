package model

import "github.com/zeromicro/go-zero/core/stores/mon"

var _ ChatFeedbackModel = (*customChatFeedbackModel)(nil)

type (
	// ChatFeedbackModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatFeedbackModel.
	ChatFeedbackModel interface {
		chatFeedbackModel
	}

	customChatFeedbackModel struct {
		*defaultChatFeedbackModel
	}
)

// NewChatFeedbackModel returns a model for the mongo.
func NewChatFeedbackModel(url, db, collection string) ChatFeedbackModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customChatFeedbackModel{
		defaultChatFeedbackModel: newDefaultChatFeedbackModel(conn),
	}
}
