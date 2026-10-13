package model

import (
	"fmt"
	"oj-server/common/proto/pb"
	"time"
)

type SubmitRecord struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

	UID         int64  `gorm:"column:uid;type:bigint;not null;index:idx_uid_created,priority:1"`
	UserName    string `gorm:"column:user_name;type:varchar(64);not null"`
	ProblemID   int64  `gorm:"column:problem_id;type:bigint;not null;index:idx_problem_created,priority:1"`
	ProblemName string `gorm:"column:problem_name;type:varchar(64);not null"`
	Lang        string `gorm:"column:lang;type:varchar(64);default:''"`
	Code        string `gorm:"column:code;type:text;not null"`
	Status      int8   `gorm:"column:status;default:0"` // 0:排队 1:判题中 2:完成（可选新增）
	Accepted    bool   `gorm:"column:accepted;default:false"`
	Message     string `gorm:"column:message;type:varchar(256);default:''"`
	Result      []byte `gorm:"column:result;type:mediumblob;not null"`
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
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UID       int64     `gorm:"column:uid;not null"`
	ProblemID int64     `gorm:"column:problem_id;not null"`
	Level     int8      `gorm:"column:level;default:0"`
	CreateAt  time.Time `gorm:"column:create_at;autoCreateTime"`
}

func (UserSolution) TableName() string {
	return "user_solution"
}

// * 在多个字段上同时声明 primaryKey 标签, GORM 会自动把它们组合成复合主键
// * 多个字段使用同一个索引名将创建复合索引，使用 priority 指定复合索引的顺序，
// 默认优先级值是 10，如果优先级值相同，则顺序取决于模型结构体字段的顺序
type Statistics struct {
	// 复合主键顺序：(Period, UID)
	// 复合索引 idx_accomplish_sort (period, accomplish_count DESC, uid)
	// 注意：DESC 索引需手写 SQL 建，GORM 标签不支持

	Period string `gorm:"column:period;type:char(7);primaryKey;not null;comment:YYYY-MM"`
	UID    int64  `gorm:"column:uid;primaryKey;not null;comment:用户id"`

	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime;comment:更新时间"`

	SubmitCount        int `gorm:"column:submit_count;default:0;comment:提交数量"`
	AccomplishCount    int `gorm:"column:accomplish_count;default:0;comment:通过数量"`
	EasyProblemCount   int `gorm:"column:easy_problem_count;default:0;comment:简单通过数"`
	MediumProblemCount int `gorm:"column:medium_problem_count;default:0;comment:中等通过数"`
	HardProblemCount   int `gorm:"column:hard_problem_count;default:0;comment:困难通过数"`
}

func (s *Statistics) TableName() string {
	return fmt.Sprintf("statistics_%s", time.Now().Format("2006"))
}
