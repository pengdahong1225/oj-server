package define

// binding的逗号之间不能有空格

// 注册登录表单
type RegisterForm struct {
	Mobile     string `form:"mobile" json:"mobile" binding:"required"`
	PassWord   string `form:"password" json:"password" binding:"required"`
	RePassWord string `form:"repassword" json:"repassword" binding:"required"`
}
type LoginFrom struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required"`
	PassWord string `form:"password" json:"password" binding:"required"`
}
type LoginWithSmsForm struct {
	Mobile     string `form:"mobile" json:"mobile" binding:"required"`
	CaptchaVal string `form:"captchaVal" json:"captchaVal" binding:"required"`
}
