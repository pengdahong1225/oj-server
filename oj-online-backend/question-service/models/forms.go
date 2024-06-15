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
