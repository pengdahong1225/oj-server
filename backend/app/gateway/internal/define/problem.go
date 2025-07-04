package define

// binding的逗号之间不能有空格

// CreateProblemForm 创建题目表单
type CreateProblemForm struct {
	Title       string   `json:"title" form:"title" binding:"required"`
	Description string   `json:"description" form:"description" binding:"required"`
	Level       int32    `json:"level" form:"level" binding:"required"`
	Tags        []string `json:"tags" form:"tags" binding:"required"`
}

// ProblemConfig 题目配置文件模型
type ProblemConfig struct {
	TestCases    []TestCase `json:"test_cases" form:"test_cases" binding:"required"`
	CompileLimit Limit      `json:"compile_limit" form:"compile_config" binding:"required"`
	RunLimit     Limit      `json:"run_limit" form:"run_config" binding:"required"`
}

type Limit struct {
	CpuLimit    int64 `json:"cpu_limit"`
	ClockLimit  int64 `json:"clock_limit"`
	MemoryLimit int64 `json:"memory_limit"`
	StackLimit  int64 `json:"stack_limit"`
	ProcLimit   int64 `json:"proc_limit"`
}

type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

// SubmitForm 代码提交表单
type SubmitForm struct {
	ProblemID int64  `json:"problem_id" form:"problem_id" binding:"required"`
	Title     string `json:"title" form:"title" binding:"required"`
	Lang      string `json:"lang" form:"lang" binding:"required"`
	Code      string `json:"code" form:"code" binding:"required"`
}

// UpdateProblemForm 修改题目表单
type UpdateProblemForm struct {
	Title  string   `json:"title" form:"title" binding:"required"`
	Level  int32    `json:"level" form:"level" binding:"required"`
	Tags   []string `json:"tags" form:"tags" binding:"required"`
	Desc   string   `json:"description" form:"description" binding:"required"`
	Config string   `json:"config" form:"config" binding:"required"`
}

// QueryProblemListParams 题目列表分页查询参数
type QueryProblemListParams struct {
	Page     int32  `json:"page" form:"page" binding:"required"`
	PageSize int32  `json:"page_size" form:"page_size" binding:"required"`
	Keyword  string `json:"keyword" form:"keyword"`
	Tag      string `json:"tag" form:"tag"`
}

// UPSSParams
// 查询题目集中哪些题目被用户 AC 了
type UPSSParams struct {
	Uid        int64   `json:"uid" form:"uid" binding:"required"`
	ProblemIds []int64 `json:"problem_ids" form:"problem_ids" binding:"required"`
}
