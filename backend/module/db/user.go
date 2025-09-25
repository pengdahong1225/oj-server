package db

import (
	"time"
)

type UserInfo struct {
	ID        int64     `gorm:"<-:false;primary_key;autoIncrement;column:id"`
	CreateAt  time.Time `gorm:"<-:false;column:create_at"`
	DeletedAt time.Time `gorm:"<-:false;column:delete_at"`
	Mobile    int64     `gorm:"column:mobile;unique"`
	PassWord  string    `gorm:"column:password"`
	NickName  string    `gorm:"column:nickname"`
	Email     string    `gorm:"column:email"`
	Gender    int32     `gorm:"column:gender"`
	Role      int32     `gorm:"column:role"`
	AvatarUrl string    `gorm:"column:avatar_url"`
}

func (UserInfo) TableName() string {
	return "user_info"
}
