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
