package model

// 题目基础信息
type Problem struct {
	ID          int64    `json:"problem_id"`
	CreateAt    int64    `json:"create_at"`
	UpdateAt    int64    `json:"update_at"`
	Title       string   `json:"problem_title"`
	Level       int32    `json:"level"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Status      int32    `json:"status"`
}
type ProblemConfig struct {
	TestCases    []TestCase `json:"test_cases"`
	CompileLimit Limit      `json:"compile_limit"`
	RunLimit     Limit      `json:"run_limit"`
}
type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}
type Limit struct {
	CpuLimit    int64 `json:"cpu_limit"`
	ClockLimit  int64 `json:"clock_limit"`
	MemoryLimit int64 `json:"memory_limit"`
	StackLimit  int64 `json:"stack_limit"`
	ProcLimit   int64 `json:"proc_limit"`
}

// 题目列表分页查询参数
type QueryProblemListRequest struct {
	Page     int32  `form:"page" binding:"required"`
	PageSize int32  `form:"page_size" binding:"required"`
	Keyword  string `form:"keyword"`
	Tag      string `form:"tag"`
}

// 创建题目表单
type CreateProblemRequest struct {
	Title string   `json:"problem_title" binding:"required"`
	Level int32    `json:"level" binding:"required"`
	Tags  []string `json:"tags" binding:"required"`
}

// 修改题目表单
type UpdateProblemRequest struct {
	ProblemID int64    `json:"problem_id" binding:"required"`
	Title     string   `json:"problem_title"  binding:"required"`
	Level     int32    `json:"level" binding:"required"`
	Tags      []string `json:"tags"  binding:"required"`
}

// 代码提交表单
type SubmitRequest struct {
	ProblemID int64  `json:"problem_id" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Lang      string `json:"lang" binding:"required"`
	Code      string `json:"code" binding:"required"`
}

// 查询题目集中哪些题目被用户 AC 了
type UPSSParams struct {
	Uid        int64   `json:"uid" form:"uid" binding:"required"`
	ProblemIds []int64 `json:"problem_ids" form:"problem_ids" binding:"required"`
}
