package model

import "oj-server/common/proto/pb"

// binding的逗号之间不能有空格

// ============================ 注册 ============================
type CreateUserRequest struct {
	Mobile   string `form:"mobile" binding:"required"`
	PassWord string `form:"password" binding:"required"`
	Username string `form:"username" binding:"required"`
	Email    string `form:"email"`
	Role     int32  `form:"role"` // 0:user 1:admin
}

// ============================ 登录 ============================
type LoginRequest struct {
	Mobile   string `form:"mobile" binding:"required"`
	PassWord string `form:"password" binding:"required"`
}
type LoginWithSmsRequest struct {
	Mobile     string `form:"mobile" binding:"required"`
	CaptchaVal string `form:"captchaVal" binding:"required"`
}
type LoginResponse struct {
	UserInfo    *pb.UserLoginResponse `json:"user_info"`
	AccessToken string                `json:"access_token"`
}

// ============================ 重置密码 ============================
type ResetPasswordRequest struct {
	Mobile     string `form:"mobile" binding:"required"`
	CaptchaVal string `form:"captchaVal" binding:"required"`
	PassWord   string `form:"password" binding:"required"`
}

// 用户列表分页查询参数
type QueryUserListRequest struct {
	Page     int32  `form:"page" binding:"required"`
	PageSize int32  `form:"page_size" binding:"required"`
	Keyword  string `form:"keyword"`
}

// 发送短信验证码
type SendSmsCodeRequest struct {
	CaptchaID   string `form:"captchaID" binding:"required"`
	CaptchaCode string `form:"captchaCode" binding:"required"`
	Mobile      string `form:"mobile" binding:"required"`
}
