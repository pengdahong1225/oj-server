package mysql

import "time"

type UserInfo struct {
	ID        int64     `gorm:"<-:false;primary_key;autoIncrement;column:id" json:"id"`
	CreateAt  time.Time `gorm:"<-:false;column:create_at" json:"createAt"`
	DeletedAt time.Time `gorm:"<-:false;column:delete_at" json:"deleteAt"`

	Mobile    int64  `gorm:"column:mobile;unique" json:"mobile"`
	PassWord  string `gorm:"column:password" json:"password"`
	NickName  string `gorm:"column:nickname" json:"userName"`
	Email     string `gorm:"column:email" json:"email"`
	Gender    int32  `gorm:"column:gender" json:"gender"`
	Role      int32  `gorm:"column:role" json:"role"`
	AvatarUrl string `gorm:"column:avatar_url" json:"avatar_url"`
}

func (UserInfo) TableName() string {
	return "user_info"
}

type UserSolution struct {
	ID        int64     `gorm:"<-:false;primary_key;autoIncrement;column:id" json:"id"`
	CreateAt  time.Time `gorm:"<-:false;column:create_at" json:"createAt"`
	DeletedAt time.Time `gorm:"<-:false;column:delete_at" json:"deleteAt"`

	Uid       int64 `gorm:"column:uid" json:"uid"`
	ProblemID int64 `gorm:"column:problem_id" json:"problem_id"`
}

func (UserSolution) TableName() string {
	return "user_solution"
}
