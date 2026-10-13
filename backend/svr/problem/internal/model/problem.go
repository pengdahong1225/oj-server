package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	// 0: 待校验 1: 未通过校验 2: 待发布 3: 已发布 4: 已下线
	ProblemStatusPending = iota
	ProblemStatusFailed
	ProblemStatusPendingPublish
	ProblemStatusPublished
	ProblemStatusOffline
)

type Problem struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	CreateAt  time.Time      `gorm:"column:create_at;autoCreateTime"`
	UpdateAt  time.Time      `gorm:"column:update_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index"`

	Title      string `gorm:"column:title;type:varchar(64);not null;uniqueIndex:idx_title"`
	Level      int8   `gorm:"column:level;default:0"`
	Tags       []byte `gorm:"column:tags;type:json"`
	PackageDir string `gorm:"column:package_dir;type:varchar(256);default:''"`
	Status     int8   `gorm:"column:status;default:0"` // 0：待校验 1：校验中 2：未通过校验 3：待发布 4：已发布 5：已下线

	SubmitCount  int64 `gorm:"column:submit_count;default:0"`
	AcceptCount  int64 `gorm:"column:accept_count;default:0"`
	CommentCount int64 `gorm:"column:comment_count;default:0"`
}

func (p *Problem) TableName() string {
	return "problem"
}
