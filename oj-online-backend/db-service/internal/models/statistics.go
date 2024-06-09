package models

import "time"

type Statistics struct {
	ID        int64     `gorm:"<-:false;primary_key;autoIncrement;column:id" json:"id"`
	CreateAt  time.Time `gorm:"<-:false;column:create_at" json:"createAt"`
	DeletedAt time.Time `gorm:"<-:false;column:delete_at" json:"deleteAt"`

	Uid                int64 `gorm:"column:uid" json:"uid"`
	SubmitCount        int64 `gorm:"column:submit_count" json:"submit_count"`
	AccomplishCount    int64 `gorm:"column:accomplish_count" json:"accomplish_count"`
	EasyProblemCount   int64 `gorm:"column:easy_problem_count" json:"easy_problem_count"`
	MediumProblemCount int64 `gorm:"column:medium_problem_count" json:"medium_problem_count"`
	HardProblemCount   int64 `gorm:"column:hard_problem_count" json:"hard_problem_count"`

	User UserInfo `gorm:"foreignKey:Uid;references:id" json:"userInfo"`
}

func (Statistics) TableName() string {
	return "user_problem_statistics"
}

type SubMitRecord struct {
	ID        int64     `gorm:"<-:false;primary_key;autoIncrement;column:id" json:"id"`
	CreateAt  time.Time `gorm:"<-:false;column:create_at" json:"createAt"`
	DeletedAt time.Time `gorm:"<-:false;column:delete_at" json:"deleteAt"`

	Uid       int64  `gorm:"column:uid" json:"uid"`
	ProblemID int64  `gorm:"column:problem_id" json:"problem_id"`
	Code      string `gorm:"column:code" json:"code"`
	Result    string `gorm:"column:result" json:"result"`
	Lang      string `gorm:"column:lang" json:"lang"`

	User    UserInfo `gorm:"foreignKey:Uid;references:id" json:"userInfo"`
	Problem Problem  `gorm:"foreignKey:ProblemID;references:id" json:"problem"`
}

func (SubMitRecord) TableName() string {
	return "user_submit_record"
}
