package model

import (
	"fmt"
	"gorm.io/gorm"
	"oj-server/pkg/proto/pb"
	"time"
)

type SubmitRecord struct {
	gorm.Model
	Uid         int64  `gorm:"column:uid;index:idx_uid" json:"uid"`
	UserName    string `gorm:"column:user_name" json:"user_name"`
	ProblemID   int64  `gorm:"column:problem_id" json:"problem_id"`
	ProblemName string `gorm:"column:problem_name" json:"problem_name"`
	Status      string `gorm:"column:status" json:"status"`
	Code        string `gorm:"column:code" json:"code"`
	Result      []byte `gorm:"column:result" json:"result"`
	Lang        string `gorm:"column:lang" json:"lang"`
}

func (receiver *SubmitRecord) TableName() string {
	return "user_submit_record"
}

func (receiver *SubmitRecord) Transform() *pb.SubmitRecord {
	return &pb.SubmitRecord{
		Id:          int64(receiver.ID),
		CreatedAt:   receiver.CreatedAt.Unix(),
		Uid:         receiver.Uid,
		UserName:    receiver.UserName,
		ProblemId:   receiver.ProblemID,
		ProblemName: receiver.ProblemName,
		Status:      receiver.Status,
		Code:        receiver.Code,
		Result:      receiver.Result,
	}
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

func (s *Statistics) TableName() string {
	// 从Period字段中提取年份
	year := ""
	if len(s.Period) >= 4 {
		year = s.Period[:4]
	} else {
		// 如果Period格式不正确，使用当前年份作为后备
		year = time.Now().Format("2006")
	}
	// 验证年份是否有效
	if _, err := time.Parse("2006", year); err != nil {
		year = time.Now().Format("2006")
	}

	return fmt.Sprintf("statistics_%s", year)
}
