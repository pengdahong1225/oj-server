package model

// 获取图形验证码返回结果
type ImageCaptchaData struct {
	CaptchaID string `json:"captcha_id"`
	Captcha   string `json:"captcha"`
}

// 获取短信验证码
type GetSmsCodeForm struct {
	CaptchaID    string `form:"captchaID" binding:"required"`
	CaptchaValue string `form:"captchaValue" binding:"required"`
	Mobile       string `form:"mobile" binding:"required"`
}
