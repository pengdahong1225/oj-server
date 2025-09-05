package db

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
	Uid       int64     `gorm:"column:uid" json:"uid"`
	ProblemID int64     `gorm:"column:problem_id" json:"problem_id"`
	CreateAt  time.Time `gorm:"<-:false;column:create_at" json:"createAt"`
	DeletedAt time.Time `gorm:"<-:false;column:delete_at" json:"deleteAt"`
}

func (UserSolution) TableName() string {
	return "user_solution"
}

type Statistics struct {
	Uid       int64     `gorm:"column:uid" json:"uid"`
	CreateAt  time.Time `gorm:"<-:false;column:create_at" json:"createAt"`
	DeletedAt time.Time `gorm:"<-:false;column:delete_at" json:"deleteAt"`

	SubmitCount        int32 `gorm:"column:submit_count" json:"submit_count"`
	AccomplishCount    int32 `gorm:"column:accomplish_count" json:"accomplish_count"`
	EasyProblemCount   int32 `gorm:"column:easy_problem_count" json:"easy_problem_count"`
	MediumProblemCount int32 `gorm:"column:medium_problem_count" json:"medium_problem_count"`
	HardProblemCount   int32 `gorm:"column:hard_problem_count" json:"hard_problem_count"`
}

func (Statistics) TableName() string {
	return "user_problem_statistics"
}
