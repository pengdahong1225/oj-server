package models

import "time"

type UserInfo struct {
	ID          int64     `gorm:"<-:false;primary_key;autoIncrement;column:id" json:"id"`
	CreateAt    time.Time `gorm:"<-:false;column:create_at" json:"createAt"`
	DeletedAt   time.Time `gorm:"<-:false;column:delete_at" json:"deleteAt"`
	Phone       int64     `gorm:"column:phone;unique" json:"phone"`
	Password    string    `gorm:"column:password" json:"password"`
	NickName    string    `gorm:"default:新用户;column:nickname" json:"userName"`
	Email       string    `gorm:"column:email" json:"email"`
	Gender      int32     `gorm:"column:gender" json:"gender"`
	Role        int32     `gorm:"column:role" json:"role"`
	HeadUrl     string    `gorm:"column:head_url" json:"headPic"`
	PassCount   int64     `gorm:"column:pass_count" json:"passCount"`
	SubmitCount int64     `gorm:"column:submit_count" json:"submitCount"`
}

func (UserInfo) TableName() string {
	return "user_info"
}

type UserSubMit struct {
	ID         int64     `gorm:"<-:false;primary_key;autoIncrement;column:id" json:"id"`
	CreateAt   time.Time `gorm:"<-:false;column:create_at" json:"createAt"`
	DeletedAt  time.Time `gorm:"<-:false;column:delete_at" json:"deleteAt"`
	UserID     int64     `gorm:"column:user_id" json:"userId"`
	QuestionID int64     `gorm:"column:question_id" json:"questionId"`
	Code       string    `gorm:"column:code" json:"code"`
	Result     string    `gorm:"column:result" json:"result"`
	Lang       string    `gorm:"column:lang" json:"lang"`
}

func (UserSubMit) TableName() string {
	return "user_submit"
}
