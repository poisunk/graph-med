package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatType struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// TODO: Fill your own fields
	TypeId    string `bson:"typeId" json:"typeId"`
	ModelName string `bson:"modelName" json:"modelName"`
	
	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
