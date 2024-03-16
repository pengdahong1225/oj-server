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

// QuestionForm 代码运行和提交表单
type QuestionForm struct {
	Id     int64  `json:"id" form:"id" binding:"required"`
	UserId int64  `json:"userId" form:"userId" binding:"required"`
	Title  string `json:"title" form:"title" binding:"required"`
	Code   string `json:"code" form:"code" binding:"required"`
	Clang  string `json:"clang" form:"clang" binding:"required"`
}

type JudgeBackForm struct {
	SessionID  string `json:"sessionID" form:"sessionID" binding:"required"`
	QuestionID int64  `json:"questionID" form:"questionID" binding:"required"`
	UserID     int64  `json:"userID" form:"userID" binding:"required"`
	Clang      string `json:"clang" form:"clang" binding:"required"`
	Status     int32  `json:"status" form:"status" binding:"required"` // 0: 正常 1: 代码非法 2: 编译错误 3: 运行超时 4: 内存溢出 5: 系统错误
	Tips       string `json:"tips" form:"tips" binding:"required"`     // 系统提示
	Output     string `json:"output" form:"output" binding:"required"` // 系统输出
}
