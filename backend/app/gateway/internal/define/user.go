package define

// binding的逗号之间不能有空格

// 注册表单
type RegisterForm struct {
	Mobile     string `form:"mobile" json:"mobile" binding:"required"`
	PassWord   string `form:"password" json:"password" binding:"required"`
	RePassWord string `form:"repassword" json:"repassword" binding:"required"`
}

// 登录表单
type LoginFrom struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required"`
	PassWord string `form:"password" json:"password" binding:"required"`
}
type LoginWithSmsForm struct {
	Mobile     string `form:"mobile" json:"mobile" binding:"required"`
	CaptchaVal string `form:"captchaVal" json:"captchaVal" binding:"required"`
}

// 重置密码表单
type ResetPasswordForm struct {
	Mobile     string `form:"mobile" json:"mobile" binding:"required"`
	CaptchaVal string `form:"captchaVal" json:"captchaVal" binding:"required"`
	PassWord   string `form:"password" json:"password" binding:"required"`
}
