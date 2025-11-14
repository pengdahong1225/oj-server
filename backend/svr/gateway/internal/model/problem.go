package model

// 题目基础信息
type Problem struct {
	ID          int64    `json:"problem_id"`
	CreateAt    string   `json:"create_at"`
	Title       string   `json:"problem_title"`
	Level       int32    `json:"level"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Status      int32    `json:"status"`
}

// 题目配置文件模型
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

// binding的逗号之间不能有空格

// 创建题目表单
type CreateProblemForm struct {
	Title       string   `json:"title" binding:"required"`
	Level       int32    `json:"level" binding:"required"`
	Tags        []string `json:"tags" binding:"required"`
	Description string   `json:"description" binding:"required"`
}

// 代码提交表单
type SubmitForm struct {
	ProblemID int64  `json:"problem_id" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Lang      string `json:"lang" binding:"required"`
	Code      string `json:"code" binding:"required"`
}
type SubmitResult struct {
	TaskId string `json:"task_id"`
}

// 修改题目表单
type UpdateProblemForm struct {
	Title  string   `json:"title" form:"title" binding:"required"`
	Level  int32    `json:"level" form:"level" binding:"required"`
	Tags   []string `json:"tags" form:"tags" binding:"required"`
	Desc   string   `json:"description" form:"description" binding:"required"`
	Config string   `json:"config" form:"config" binding:"required"`
}

// 题目列表分页查询参数
type QueryProblemListParams struct {
	Page     int32  `form:"page" binding:"required"`
	PageSize int32  `form:"page_size" binding:"required"`
	Keyword  string `form:"keyword"`
	Tag      string `form:"tag"`
}
type QueryProblemListResult struct {
	Total int64      `json:"total"`
	List  []*Problem `json:"list"`
}

// 查询题目集中哪些题目被用户 AC 了
type UPSSParams struct {
	Uid        int64   `json:"uid" form:"uid" binding:"required"`
	ProblemIds []int64 `json:"problem_ids" form:"problem_ids" binding:"required"`
}
