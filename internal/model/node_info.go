package model

import "time"

type NodeInfo struct {
	ID        int       `xorm:"pk autoincr 'id'"`
	Label     string    `xorm:"varchar(100) not null 'label'"`
	Name      string    `xorm:"varchar(100) not null 'name'"`
	AttrName  string    `xorm:"varchar(100) not null 'attr_name'"`
	AttrValue string    `xorm:"longtext not null 'attr_value'"`
	Group     string    `xorm:"varchar(100) not null 'group'"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}
