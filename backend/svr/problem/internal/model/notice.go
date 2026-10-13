package model

import (
	"time"

	"oj-server/common/proto/pb"

	"gorm.io/gorm"
)

type Notice struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	CreateAt  time.Time      `gorm:"column:create_at;autoCreateTime"`
	UpdateAt  time.Time      `gorm:"column:update_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index"`

	Title    string `gorm:"column:title;type:varchar(64);not null"`
	Content  string `gorm:"column:content;type:text;not null"`
	CreateBy int64  `gorm:"column:create_by;default:0"`
	Status   int8   `gorm:"column:status;default:0"` // 1：正常 0：隐藏
}

func (n *Notice) TableName() string {
	return "notice"
}

// Transform2Pb 转换为 protobuf 结构
func (n *Notice) Transform2Pb() *pb.Notice {
	return &pb.Notice{
		Id:       n.ID,
		Title:    n.Title,
		CreateAt: n.CreateAt.Format("2006-01-02 15:04:05"),
		Content:  n.Content,
		Status:   int32(n.Status),
		CreateBy: n.CreateBy,
	}
}
