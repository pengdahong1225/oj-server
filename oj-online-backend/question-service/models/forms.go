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

// QuestionForm 代码运行和提交表单
type QuestionForm struct {
	Id     int64  `json:"id" form:"id" binding:"required"`
	UserId int64  `json:"userId" form:"userId" binding:"required"`
	Title  string `json:"title" form:"title" binding:"required"`
	Code   string `json:"code" form:"code" binding:"required"`
	Clang  string `json:"clang" form:"clang" binding:"required"`
}
