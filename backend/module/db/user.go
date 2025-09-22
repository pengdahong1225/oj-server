package db

import (
	"fmt"
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

type UserSolution struct {
	ID        int64     `gorm:"<-:false;primary_key;autoIncrement;column:id"`
	Uid       int64     `gorm:"column:uid"`
	ProblemID int64     `gorm:"column:problem_id"`
	CreateAt  time.Time `gorm:"<-:false;column:create_at"`
	DeletedAt time.Time `gorm:"<-:false;column:delete_at"`
}

func (UserSolution) TableName() string {
	return "user_solution"
}

type Statistics struct {
	Uid                int64     `gorm:"column:uid;primaryKey"`
	Period             string    `gorm:"column:period;primaryKey;size:7;index:idx_period"`
	CreateAt           time.Time `gorm:"column:create_at;autoCreateTime;index:idx_create_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime"`
	SubmitCount        int32     `gorm:"column:submit_count;default:0"`
	AccomplishCount    int32     `gorm:"column:accomplish_count;default:0;index:idx_accomplish"`
	EasyProblemCount   int32     `gorm:"column:easy_problem_count;default:0"`
	MediumProblemCount int32     `gorm:"column:medium_problem_count;default:0"`
	HardProblemCount   int32     `gorm:"column:hard_problem_count;default:0"`
}

//func (Statistics) TableName() string {
//	return "user_problem_statistics"
//}

func (Statistics) TableName(stamp int64) string {
	return fmt.Sprintf("user_submit_record_%d", time.Unix(stamp, 0).Year())
}
