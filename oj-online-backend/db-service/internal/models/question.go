package models

import "time"

type Question struct {
	ID        int64     `gorm:"<-:false;primary_key;autoIncrement;column:id" json:"id"`
	CreateAt  time.Time `gorm:"<-:false;column:create_at" json:"createAt"`
	DeletedAt time.Time `gorm:"<-:false;column:delete_at" json:"deleteAt"`

	Title string `gorm:"column:title" json:"title"`
	Level int32  `gorm:"column:level" json:"level"`
	Tags  string `gorm:"column:tags" json:"tags"`

	Description string `gorm:"column:description" json:"description"` // markdown
	TestCase    string `gorm:"column:test_case" json:"test_case"`     // json{input,output}
	Template    string `gorm:"column:template" json:"template"`       // json{lang,code}
}

func (Question) TableName() string {
	return "question"
}