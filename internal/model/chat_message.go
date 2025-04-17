package model

import "time"

type ChatMessage struct {
	ID              int       `xorm:"pk autoincr 'id'"`
	MessageID       int       `xorm:"int not null 'message_id'"`
	ParentMessageID int       `xorm:"int not null 'parent_message_id'"`
	UserID          string    `xorm:"varchar(64) not null 'user_id'"`
	ChatSessionID   string    `xorm:"varchar(64) not null 'chat_session_id'"`
	ModelName       string    `xorm:"varchar(100) not null"`
	Role            string    `xorm:"enum('user', 'assistant', 'system') not null"`
	Content         string    `xorm:"longtext"`
	Operator        string    `xorm:"varchar(64)"`
	CreatedAt       time.Time `xorm:"created"`
	UpdatedAt       time.Time `xorm:"updated"`
	DeletedAt       time.Time `xorm:"deleted"`
}

func (ChatMessage) TableName() string {
	return "chat_message"
}
