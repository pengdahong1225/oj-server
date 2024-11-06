package mysql

import (
	"gorm.io/gorm"
)

type SubmitRecord struct {
	gorm.Model
	Uid         int64  `gorm:"column:uid;index:idx_uid" json:"uid"`
	ProblemID   int64  `gorm:"column:problem_id" json:"problem_id"`
	ProblemName string `gorm:"column:problem_name" json:"problem_name"`
	Status      int32  `gorm:"column:status" json:"status"`
	Code        string `gorm:"column:code" json:"code"`
	Result      []byte `gorm:"column:result type:blob" json:"result"`
	Lang        string `gorm:"column:lang" json:"lang"`
}

func (receiver SubmitRecord) TableName() string {
	return "user_submit_record"
}

//func (receiver *SubmitRecord) TableName(stamp int64) string {
//	t := time.Unix(stamp, 0)
//	s := t.Format("20060102")
//	return fmt.Sprintf("user_submit_record_%s", s)
//}
