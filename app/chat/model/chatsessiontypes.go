package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatSession struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// TODO: Fill your own fields
	SessionId string    `bson:"sessionId" json:"sessionId"`
	TypeId    string    `bson:"typeId" json:"typeId"`
	UserId    string    `bson:"userId" json:"userId"`
	UpdateAt  time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt  time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
