package model

import "time"

type ChatFeedback struct {
	ID            int       `xorm:"pk autoincr 'id'"`
	MessageID     int       `xorm:"int not null 'message_id'"`
	UserID        string    `xorm:"varchar(64) not null 'user_id'"`
	ChatSessionID string    `xorm:"varchar(64) not null 'chat_session_id'"`
	Feedback      string    `xorm:"varchar(255)"`
	CreatedAt     time.Time `xorm:"created"`
	UpdatedAt     time.Time `xorm:"updated"`
	DeletedAt     time.Time `xorm:"deleted"`
}

func (ChatFeedback) TableName() string {
	return "chat_feedback"
}
