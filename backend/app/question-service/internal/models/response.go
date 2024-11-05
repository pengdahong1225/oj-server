package models

import "github.com/pengdahong1225/oj-server/backend/proto/pb"

// Response 统一返回格式
type Response struct {
	Code    int    `json:"code"` // 业务状态码
	Message string `json:"message"`
	Data    any    `json:"data"`
}

const (
	Success = 0
	Failed  = 1
)

type RankList struct {
	Phone     int64  `json:"phone"`
	NickName  string `json:"nickName"`
	PassCount int64  `json:"passCount"`
}

type LoginRspData struct {
	Uid       int64  `json:"uid"`
	Mobile    int64  `json:"mobile"`
	NickName  string `json:"nickname"`
	Email     string `json:"email"`
	Gender    int32  `json:"gender"`
	Role      int32  `json:"role"`
	AvatarUrl string `json:"avatar_url"`
	Token     string `json:"token"`
}

type NoticeRspData struct {
	Total      int32        `json:"total"`
	NoticeList []*pb.Notice `json:"noticeList"`
}

// QueryRecordResponse 查询提交记录响应格式
// 根据 uid problemID stamp 查询用户在某天的对某题目的提交记录
type QueryRecordResponse struct {
	Uid       int64 `json:"uid"`
	ProblemID int64 `json:"problem_id"`
	Stamp     int64 `json:"stamp"`

	Records []SubmitRecord `json:"records"` // 提交记录列表
}
type SubmitRecord struct {
	Code   string            `json:"code"`
	Lang   string            `json:"lang"`
	Result []*pb.JudgeResult `json:"result"` // 判题结果集
}

// QuerySubmitResultResponse 查询提交结果响应格式
type QuerySubmitResultResponse struct {
	Uid       int64             `json:"uid"`
	ProblemID int64             `json:"problem_id"`
	Result    []*pb.JudgeResult `json:"result"` // 判题结果集
}
