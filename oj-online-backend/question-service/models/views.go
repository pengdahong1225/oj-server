package models

import "github.com/golang/protobuf/ptypes/timestamp"

type RankList struct {
	Phone     int64  `json:"phone"`
	NickName  string `json:"nickName"`
	PassCount int64  `json:"passCount"`
}

type UserInfo struct {
	CreateAt    *timestamp.Timestamp `json:"createAt"`
	Phone       int64                `json:"phone"`
	NickName    string               `json:"nickname"`
	Email       string               `json:"email"`
	Gender      int32                `json:"gender"`
	Role        int32                `json:"role"`
	HeadUrl     string               `json:"head_url"`
	PassCount   int64                `json:"passCount"`
	SubmitCount int64                `json:"submitCount"`
}

type Question struct {
	Id          int64                `json:"id"`
	CreateAt    *timestamp.Timestamp `json:"createAt"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Level       int32                `json:"level"`
	Tags        []string             `json:"tags"`
}

type QuestionResult struct {
	QuestionID int64  `json:"questionID"`
	UserID     int64  `json:"userID"`
	Clang      string `json:"clang"`
	Status     int32  `json:"status"` // 0: 正常 1: 代码非法 2: 编译错误 3: 运行超时 4: 内存溢出 5: 系统错误
	Tips       string `json:"tips"`   // 系统提示
	Output     string `json:"output"` // 系统输出
}
