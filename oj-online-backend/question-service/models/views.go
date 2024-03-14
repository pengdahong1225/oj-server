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
	Id     int64  `json:"id" form:"id" binding:"required"`
	UserId int64  `json:"userId" form:"userId" binding:"required"`
	Clang  string `json:"clang" form:"clang" binding:"required"`
	Result string `json:"result"`
}
