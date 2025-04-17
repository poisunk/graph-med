package model

import "time"

type ChatType struct {
	ID        int       `xorm:"pk autoincr 'id'"`
	TypeID    string    `xorm:"varchar(64) not null 'type_id'"`
	TypeName  string    `xorm:"varchar(64) not null 'type_name'"`
	McpIDs    string    `xorm:"varchar(255) 'mcp_ids'"`
	ModelName string    `xorm:"varchar(64) 'model_name'"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

func (ChatType) TableName() string {
	return "chat_type"
}
