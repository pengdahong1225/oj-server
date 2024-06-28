package models

type TestCase struct {
	Id int32 `json:"id"`
}

// 用户状态枚举
const (
	UserStateNormal = iota
	UserStateJudging
)

// 用户提交题目状态枚举
const (
	UPStateNormal    = iota
	UPStateCompiling // 编译中
	UPStateRunning   // 运行中
	UPStateJudging   // 判题中
)
