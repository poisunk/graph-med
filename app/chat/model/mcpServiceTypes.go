package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type McpService struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// TODO: Fill your own fields
	McpId   string `bson:"mcpId" json:"mcpId"`
	Type    string `bson:"type" json:"type"`
	Name    string `bson:"name" json:"name"`
	Command string `bson:"command" json:"command"`
	BaseUrl string `bson:"baseUrl" json:"baseUrl"`

	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
