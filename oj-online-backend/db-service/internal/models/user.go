package models

import "time"

type UserInfo struct {
	ID        int64     `gorm:"<-:false;primary_key;autoIncrement;column:id" json:"id"`
	CreateAt  time.Time `gorm:"<-:false;column:create_at" json:"createAt"`
	DeletedAt time.Time `gorm:"<-:false;column:delete_at" json:"deleteAt"`
	Phone     int64     `gorm:"column:phone;unique" json:"phone"`
	Password  string    `gorm:"column:password" json:"password"`
	NickName  string    `gorm:"default:新用户;column:nickname" json:"userName"`
	Email     string    `gorm:"column:email" json:"email"`
	Gender    int32     `gorm:"column:gender" json:"gender"`
	Role      int32     `gorm:"column:role" json:"role"`
	HeadUrl   string    `gorm:"column:head_url" json:"headPic"`
}

func (UserInfo) TableName() string {
	return "user_info"
}
