package mysql

import "time"

type Problem struct {
	ID        int64     `gorm:"<-:false;primary_key;autoIncrement;column:id" json:"id"`
	CreateAt  time.Time `gorm:"<-:false;column:create_at" json:"createAt"`
	DeletedAt time.Time `gorm:"<-:false;column:delete_at" json:"deleteAt"`

	Title        string `gorm:"column:title" json:"title"`
	Level        int32  `gorm:"column:level" json:"level"`
	Tags         string `gorm:"column:tags" json:"tags"`
	Description  string `gorm:"column:description" json:"description"`
	CreateBy     int64  `gorm:"column:create_by" json:"create_by"`
	CommentCount int64  `gorm:"column:comment_count" json:"comment_count"`

	TestCase      string `gorm:"column:test_case" json:"test_case"`
	CompileConfig string `gorm:"column:compile_config" json:"compile_config"`
	RunConfig     string `gorm:"column:run_config" json:"run_config"`
}

func (Problem) TableName() string {
	return "problem"
}

type ProblemConfig struct {
	CpuLimit    int64 `json:"cpu_limit"`
	ClockLimit  int64 `json:"clock_limit"`
	MemoryLimit int64 `json:"memory_limit"`
	ProcLimit   int64 `json:"proc_limit"`
}

type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type ProblemHotData struct {
	TestCase      string `json:"test_case"`
	CompileConfig string `json:"compile_config"`
	RunConfig     string `json:"run_config"`
}
