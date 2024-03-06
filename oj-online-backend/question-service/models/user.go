package models

// binding的逗号之间不能有空格

type RegistryForm struct {
	Phone    string `form:"phone" json:"phone" binding:"required,phone"` // 需要自定义validator
	PassWord string `form:"password" json:"password" binding:"required,min=4,max=20"`
	SmsCode  string `form:"sms_code" json:"sms_code" binding:"required"`
	NickName string `form:"nickname" json:"nickname"`
	Email    string `form:"email" json:"email"`
	Gender   int    `form:"gender" json:"gender"`
	Role     int    `form:"role" json:"role"`
	HeadUrl  string `form:"head_url" json:"head_url"`
}

type LoginFrom struct {
	Phone    string `form:"phone" json:"phone" binding:"required,phone"` // 需要自定义validator
	PassWord string `form:"password" json:"password" binding:"required,min=4,max=20"`
}
