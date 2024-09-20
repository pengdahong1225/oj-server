package models

// binding的逗号之间不能有空格

// GetSmsCodeForm 获取图像验证码表单
type GetSmsCodeForm struct {
	CaptchaID    string `form:"captchaID" json:"captchaID" binding:"required"`
	CaptchaValue string `form:"captchaValue" json:"captchaValue" binding:"required"`
	Mobile       string `form:"mobile" json:"mobile" binding:"required"`
}

type LoginFrom struct {
	Mobile  string `form:"mobile" json:"mobile" binding:"required"`
	SmsCode string `form:"smsCode" json:"smsCode" binding:"required"`
}

// SubmitForm 代码提交表单
type SubmitForm struct {
	ProblemID int64  `json:"problem_id" form:"problem_id" binding:"required"`
	Title     string `json:"title" form:"title" binding:"required"`
	Lang      string `json:"lang" form:"lang" binding:"required"`
	Code      string `json:"code" form:"code" binding:"required"`
}

// AddProblemForm 添加、修改题目表单
type AddProblemForm struct {
	Title         string        `json:"title" form:"title" binding:"required"`
	Level         int32         `json:"level" form:"level" binding:"required"`
	Tags          []string      `json:"tags" form:"tags" binding:"required"`
	Desc          string        `json:"description" form:"description" binding:"required"`
	TestCases     []TestCase    `json:"testCases" form:"testCases" binding:"required"`
	CompileConfig ProblemConfig `json:"compile_config" form:"compile_config" binding:"required"`
	RunConfig     ProblemConfig `json:"run_config" form:"run_config" binding:"required"`
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

type AddCommentForm struct {
	ObjId         int64  `json:"obj_id" form:"obj_id" binding:"required"`
	UserId        int64  `json:"user_id" form:"user_id" binding:"required"`
	UserName      string `json:"user_name" form:"user_name"`
	UserAvatarUrl string `json:"user_avatar_url" form:"user_avatar_url"`
	Content       string `json:"content" form:"content" binding:"required"`
	Stamp         int64  `json:"stamp" form:"stamp"`

	RootId         int64 `json:"root_id" form:"root_id" binding:"required"`
	RootCommentId  int64 `json:"root_comment_id" form:"root_comment_id" binding:"required"`
	ReplyId        int64 `json:"reply_id" form:"reply_id" binding:"required"`
	ReplyCommentId int64 `json:"reply_comment_id" form:"reply_comment_id" binding:"required"`
}
type QueryCommentForm struct {
	ObjId int64 `json:"obj_id" form:"obj_id" binding:"required"`

	RootId        int64 `json:"root_id" form:"root_id"`
	RootCommentId int64 `json:"root_comment_id" form:"root_comment_id"`

	ReplyId        int64 `json:"reply_id" form:"reply_id"`
	ReplyCommentId int64 `json:"reply_comment_id" form:"reply_comment_id"`

	CurSor int64 `json:"cur_cursor" form:"cur_cursor"`
}
