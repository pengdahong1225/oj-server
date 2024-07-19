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

// AddProblemForm 添加题目表单
type AddProblemForm struct {
	Title       string   `json:"title" form:"title" binding:"required"`
	Level       int32    `json:"level" form:"level" binding:"required"`
	Tags        []string `json:"tags" form:"tags" binding:"required"`
	Desc        string   `json:"description" form:"description" binding:"required"`
	TestCase    string   `json:"testCase" form:"testCase" binding:"required"`
	CpuLimit    int64    `json:"cpuLimit" form:"cpuLimit" binding:"required"`
	ClockLimit  int64    `json:"clockLimit" form:"clockLimit" binding:"required"`
	MemoryLimit int64    `json:"memoryLimit" form:"memoryLimit" binding:"required"`
	ProcLimit   int64    `json:"procLimit" form:"procLimit" binding:"required"`
}
