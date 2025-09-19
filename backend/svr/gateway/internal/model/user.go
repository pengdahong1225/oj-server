package model

import "oj-server/proto/pb"

// binding的逗号之间不能有空格

// ============================ 注册 ============================
type RegisterForm struct {
	Mobile     string `form:"mobile" binding:"required"`
	PassWord   string `form:"password" binding:"required"`
	RePassWord string `form:"repassword"  binding:"required"`
}

// ============================ 登录 ============================
type LoginFrom struct {
	Mobile   string `form:"mobile" binding:"required"`
	PassWord string `form:"password" binding:"required"`
}
type LoginWithSmsForm struct {
	Mobile     string `form:"mobile" binding:"required"`
	CaptchaVal string `form:"captchaVal" binding:"required"`
}
type LoginResponse struct {
	UserInfo    *pb.UserLoginResponse `json:"user_info"`
	AccessToken string                `json:"access_token"`
}

// ============================ 重置密码 ============================
type ResetPasswordForm struct {
	Mobile     string `form:"mobile" binding:"required"`
	CaptchaVal string `form:"captchaVal" binding:"required"`
	PassWord   string `form:"password" binding:"required"`
}
