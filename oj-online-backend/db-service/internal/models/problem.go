package models

import "time"

type Problem struct {
	ID        int64     `gorm:"<-:false;primary_key;autoIncrement;column:id" json:"id"`
	CreateAt  time.Time `gorm:"<-:false;column:create_at" json:"createAt"`
	DeletedAt time.Time `gorm:"<-:false;column:delete_at" json:"deleteAt"`

	Title string `gorm:"column:title" json:"title"`
	Level int32  `gorm:"column:level" json:"level"`
	Tags  string `gorm:"column:tags" json:"tags"`

	Description string `gorm:"column:description" json:"description"`
	TestCase    string `gorm:"column:test_case" json:"test_case"`

	TimeLimit   int32  `gorm:"column:time_limit" json:"time_limit"`
	MemoryLimit int32  `gorm:"column:memory_limit" json:"memory_limit"`
	IoMode      string `gorm:"column:io_mode" json:"io_mode"`
	CreateBy    int64  `gorm:"column:create_by" json:"create_by"`
}

func (Problem) TableName() string {
	return "problem"
}

type ProblemDataResult struct {
	Problem
	CreateUserNickName string `gorm:"column:nickname"` // 用于Scan(挂载)user_info表中的nickname字段
}
