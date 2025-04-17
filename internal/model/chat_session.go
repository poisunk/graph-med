package model

import "time"

type ChatSession struct {
	ID        int       `xorm:"pk autoincr 'id'"`
	TypeID    string    `xorm:"varchar(64) not null default 'default' 'type_id'"`
	SessionID string    `xorm:"varchar(64) not null 'session_id'"`
	UserID    string    `xorm:"varchar(64) not null 'user_id'"`
	Title     string    `xorm:"varchar(100) not null"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}
