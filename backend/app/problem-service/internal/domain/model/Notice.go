package model

import "time"

type Notice struct {
	ID        int64     `gorm:"<-:false;primary_key;autoIncrement;column:id" json:"id"`
	CreateAt  time.Time `gorm:"<-:false;column:create_at" json:"createAt"`
	UpdateAt  time.Time `gorm:"<-:false;column:update_at" json:"updateAt"`
	DeletedAt time.Time `gorm:"<-:false;column:delete_at" json:"deleteAt"`

	Title    string `gorm:"column:title" json:"title"`
	Content  string `gorm:"column:content" json:"content"`
	CreateBy int64  `gorm:"column:create_by" json:"create_by"`
	Status   int32  `gorm:"column:status" json:"status"`
}

func (Notice) TableName() string {
	return "notice"
}
