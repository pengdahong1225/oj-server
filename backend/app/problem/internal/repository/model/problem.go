package model

import (
	"github.com/pengdahong1225/oj-server/backend/module/utils"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"time"
)

type Problem struct {
	ID        int64     `gorm:"<-:false;primary_key;autoIncrement;column:id" json:"id"`
	CreateAt  time.Time `gorm:"<-:false;column:create_at" json:"createAt"`
	DeletedAt time.Time `gorm:"<-:false;column:delete_at" json:"deleteAt"`

	Title        string `gorm:"column:title" json:"title"`
	Level        int32  `gorm:"column:level" json:"level"`
	Tags         []byte `gorm:"column:tags;type:json" json:"tags"`
	Description  string `gorm:"column:description" json:"description"`
	CreateBy     int64  `gorm:"column:create_by" json:"create_by"`
	CommentCount int64  `gorm:"column:comment_count;type:blob" json:"comment_count"`

	Config []byte `gorm:"column:config" json:"config"`
}

func (p *Problem) TableName() string {
	return "problem"
}

func (p *Problem) Transform() *pb.Problem {
	return &pb.Problem{
		Id:          p.ID,
		CreateAt:    p.CreateAt.String(),
		Title:       p.Title,
		Description: p.Description,
		Level:       p.Level,
		Tags:        utils.SplitStringWithX(string(p.Tags), "#"),
		CreateBy:    p.CreateBy,
	}
}
func TransformList(list []Problem) []*pb.Problem {
	var problems []*pb.Problem
	for _, problem := range list {
		problems = append(problems, problem.Transform())
	}
	return problems
}
