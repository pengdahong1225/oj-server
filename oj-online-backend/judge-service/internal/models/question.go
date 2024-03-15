package models

// QuestionForm 代码运行和提交表单
type QuestionForm struct {
	Id     int64  `json:"id" form:"id" binding:"required"`
	UserId int64  `json:"userId" form:"userId" binding:"required"`
	Title  string `json:"title" form:"title" binding:"required"`
	Code   string `json:"code" form:"code" binding:"required"`
	Clang  string `json:"clang" form:"clang" binding:"required"`
}
