package model

import (
	"time"
)

type Problem struct {
	ID       int64      `gorm:"column:id;primaryKey;autoIncrement"`
	CreateAt time.Time  `gorm:"column:create_at;autoCreateTime"`
	UpdateAt time.Time  `gorm:"column:update_at;autoUpdateTime"`
	DeleteAt *time.Time `gorm:"column:delete_at"`

	Title        string `gorm:"column:title;type:varchar(64);not null;uniqueIndex:idx_title"`
	Level        int8   `gorm:"column:level;default:0"`
	Tags         []byte `gorm:"column:tags"`
	Description  string `gorm:"column:description;type:text;not null"`
	CommentCount int64  `gorm:"column:comment_count;default:0"`

	Status    int8   `gorm:"column:status;default:0"`
	ConfigURL string `gorm:"column:config_url;type:varchar(256);default:''"`
}

func (p *Problem) TableName() string {
	return "problem"
}
