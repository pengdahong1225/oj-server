package mysql

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
