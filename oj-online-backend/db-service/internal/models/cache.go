package models

type ProblemHotData struct {
	TestCase string `gorm:"column:test_case" json:"test_case"`

	CpuLimit    int64 `gorm:"column:cpu_limit" json:"cpu_limit"`
	ClockLimit  int64 `gorm:"column:clock_limit" json:"clock_limit"`
	TimeLimit   int64 `gorm:"column:time_limit" json:"time_limit"`
	MemoryLimit int64 `gorm:"column:memory_limit" json:"memory_limit"`
	ProcLimit   int64 `gorm:"column:proc_limit" json:"proc_limit"`
}
