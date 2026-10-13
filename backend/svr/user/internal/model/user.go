package model

import (
	"time"

	"gorm.io/gorm"
)

type UserInfo struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	CreateAt  time.Time      `gorm:"column:create_at;autoCreateTime"`
	UpdateAt  time.Time      `gorm:"column:update_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index"`

	Mobile    int64  `gorm:"column:mobile;uniqueIndex;not null"`
	Password  string `gorm:"column:password;not null"`
	Username  string `gorm:"column:username;default:新用户"`
	Email     string `gorm:"column:email;default:''"`
	Role      int32  `gorm:"column:role;default:0"`
	AvatarURL string `gorm:"column:avatar_url;default:''"`
}

func (UserInfo) TableName() string {
	return "user_info"
}
