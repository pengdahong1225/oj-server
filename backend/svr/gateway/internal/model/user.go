package model

import "oj-server/pkg/proto/pb"

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

type UserInfo struct {
	Uid       int64  `json:"uid"`
	CreateAt  int64  `json:"create_at"`
	Mobile    int64  `json:"mobile"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Gender    int32  `json:"gender"`
	Role      int32  `json:"role"`
	AvatarUrl string `json:"avatar_url"`
}

// 用户列表分页查询参数
type QueryUserListParams struct {
	Page     int32  `form:"page" binding:"required"`
	PageSize int32  `form:"page_size" binding:"required"`
	Keyword  string `form:"keyword"`
}
type QueryUserListResult struct {
	Total int64       `json:"total"`
	List  []*UserInfo `json:"list"`
}
