package models

import "github.com/golang/protobuf/ptypes/timestamp"

// Response 统一返回格式
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

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

type SubmitResponse struct {
	QuestionID int64          `json:"questionID"`
	UserID     int64          `json:"userID"`
	Clang      string         `json:"clang"`
	ResultList []SubmitResult `json:"resultList"`
}

type SubmitResult struct {
	Result   int32  `json:"result,omitempty"`
	CpuTime  int32  `json:"cpu_time,omitempty"`
	RealTime int32  `json:"real_time,omitempty"`
	Memory   int32  `json:"memory,omitempty"`
	Signal   int32  `json:"signal,omitempty"`
	ExitCode int32  `json:"exit_code,omitempty"`
	Error    int32  `json:"error,omitempty"`
	Content  string `json:"content,omitempty"`
}
