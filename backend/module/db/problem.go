package db

import (
	"time"
)

type Problem struct {
	ID        int64     `gorm:"<-:false;primary_key;autoIncrement;column:id"`
	CreateAt  time.Time `gorm:"<-:false;column:create_at"`
	DeletedAt time.Time `gorm:"<-:false;column:delete_at"`

	Title        string `gorm:"column:title"`
	Level        int32  `gorm:"column:level"`
	Tags         []byte `gorm:"column:tags;type:json"`
	Description  string `gorm:"column:description"`
	CreateBy     int64  `gorm:"column:create_by"`
	CommentCount int64  `gorm:"column:comment_count;type:blob"`
	Status       int32  `gorm:"column:status"` // 状态 1：发布 0：隐藏（默认值）
}

func (p *Problem) TableName() string {
	return "problem"
}
