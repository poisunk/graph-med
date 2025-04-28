package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatMessage struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// TODO: Fill your own fields
	SessionId       string         `bson:"sessionId" json:"sessionId"`
	MessageId       int64          `bson:"mssageId" json:"mssageId"`
	ParentMessageId int64          `bson:"parentMessageId" json:"parentMessageId"`
	Role            string         `bson:"role" json:"role"`
	Content         string         `bson:"content" json:"content"`
	Turns           []*MessageTurn `bson:"turns" json:"turns"`
	UpdateAt        time.Time      `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt        time.Time      `bson:"createAt,omitempty" json:"createAt,omitempty"`
}

type MessageTurn struct {
	Content      string `bson:"content" json:"content"`
	Type         string `bson:"type" json:"type"`
	FunctionName string `bson:"functionName" json:"functionName"`
	FunctionArgs string `bson:"functionArgs" json:"functionArgs"`
}
