package model

import "time"

type ChatUser struct {
	ID        int       `xorm:"pk autoincr 'id'"`
	UserID    string    `xorm:"varchar(64) not null 'user_id'"`
	Username  string    `xorm:"varchar(64) not null 'username'"`
	Password  string    `xorm:"varchar(64) not null 'password'"`
	Email     string    `xorm:"varchar(64) not null 'email'"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

func (ChatUser) TableName() string {
	return "chat_user"
}
