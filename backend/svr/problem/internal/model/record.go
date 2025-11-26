package model

import (
	"fmt"
	"oj-server/pkg/proto/pb"
	"time"
)

type SubmitRecord struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	UID         int64  `gorm:"type:BIGINT;not null;comment:用户id"`
	UserName    string `gorm:"type:varchar(64);not null;comment:用户名"`
	ProblemID   int64  `gorm:"type:BIGINT;not null;comment:题目id"`
	ProblemName string `gorm:"type:varchar(64);not null;comment:题目名称"`

	Lang string `gorm:"type:varchar(64);default:'';comment:语言"`
	Code string `gorm:"type:text;not null;comment:提交的代码"`

	Accepted bool   `gorm:"default:false;comment:是否通过"`
	Message  string `gorm:"type:varchar(256);default:'';comment:测评结果描述"`
	Result   []byte `gorm:"type:blob;not null;comment:运行结果集"`
}

func (receiver *SubmitRecord) TableName() string {
	return "user_submit_record"
}

func (receiver *SubmitRecord) Transform() *pb.SubmitRecord {
	return &pb.SubmitRecord{
		Id:        int64(receiver.ID),
		CreatedAt: receiver.CreatedAt.Unix(),
		Uid:       receiver.UID,
		ProblemId: receiver.ProblemID,
		Message:   receiver.Message,
		Code:      receiver.Code,
		Results:   receiver.Result,
		Lang:      receiver.Lang,
		Accepted:  receiver.Accepted,
	}
}

type UserSolution struct {
	ID        int64     `gorm:"<-:false;primary_key;autoIncrement;column:id"`
	Uid       int64     `gorm:"column:uid"`
	ProblemID int64     `gorm:"column:problem_id"`
	CreateAt  time.Time `gorm:"<-:false;column:create_at"`
}

func (UserSolution) TableName() string {
	return "user_solution"
}

// * 在多个字段上同时声明 primaryKey 标签, GORM 会自动把它们组合成复合主键
// * 多个字段使用同一个索引名将创建复合索引，使用 priority 指定复合索引的顺序，
// 默认优先级值是 10，如果优先级值相同，则顺序取决于模型结构体字段的顺序
type Statistics struct {
	// 复合主键：(UID, Period)
	// 复合索引: idx_accomplish_sort,sort:desc (period, accomplish_count DESC, uid)

	UID       int64     `gorm:"column:uid;primaryKey;not null;comment:用户id;index:idx_accomplish_sort,sort:desc,priority:3"`
	Period    string    `gorm:"column:period;type:char(7);primaryKey;index:idx_accomplish_sort,sort:desc,priority:1;not null;comment:YYYY-MM"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime;comment:更新时间"`

	SubmitCount        int `gorm:"column:submit_count;default:0;comment:提交数量"`
	AccomplishCount    int `gorm:"column:accomplish_count;default:0;comment:通过数量;index:idx_accomplish_sort,sort:desc,priority:2;"`
	EasyProblemCount   int `gorm:"column:easy_problem_count;default:0;comment:简单通过数"`
	MediumProblemCount int `gorm:"column:medium_problem_count;default:0;comment:中等通过数"`
	HardProblemCount   int `gorm:"column:hard_problem_count;default:0;comment:困难通过数"`
}

func (s *Statistics) TableName() string {
	return fmt.Sprintf("statistics_%s", time.Now().Format("2006"))
}
