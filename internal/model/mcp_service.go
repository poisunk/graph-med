package model

import "time"

type McpService struct {
	ID        int       `xorm:"pk autoincr 'id'"`
	Name      string    `xorm:"varchar(64) not null 'name'"`
	McpID     string    `xorm:"varchar(64) not null 'mcp_id'"`
	Type      string    `xorm:"varchar(64) not null 'type'"`
	Args      string    `xorm:"longtext 'args'"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

func (McpService) TableName() string {
	return "mcp_service"
}
