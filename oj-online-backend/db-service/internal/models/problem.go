package models

import "time"

type Problem struct {
	ID        int64     `gorm:"<-:false;primary_key;autoIncrement;column:id" json:"id"`
	CreateAt  time.Time `gorm:"<-:false;column:create_at" json:"createAt"`
	DeletedAt time.Time `gorm:"<-:false;column:delete_at" json:"deleteAt"`

	Title       string `gorm:"column:title" json:"title"`
	Level       int32  `gorm:"column:level" json:"level"`
	Tags        string `gorm:"column:tags" json:"tags"`
	Description string `gorm:"column:description" json:"description"`
	TestCase    string `gorm:"column:test_case" json:"test_case"`

	CpuLimit    int64 `gorm:"column:cpu_limit" json:"cpu_limit"`
	ClockLimit  int64 `gorm:"column:clock_limit" json:"clock_limit"`
	TimeLimit   int64 `gorm:"column:time_limit" json:"time_limit"`
	MemoryLimit int64 `gorm:"column:memory_limit" json:"memory_limit"`
	ProcLimit   int64 `gorm:"column:proc_limit" json:"proc_limit"`

	CreateBy int64 `gorm:"column:create_by" json:"create_by"`
}

func (Problem) TableName() string {
	return "problem"
}
