package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatFeedback struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// TODO: Fill your own fields
	UserId    string `bson:"userId" json:"userId"`
	SessionId string `bson:"sessionId" json:"sessionId"`
	MessageId int64  `bson:"messageId" json:"messageId"`
	Feedback  string `bson:"feedback" json:"feedback"`

	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
