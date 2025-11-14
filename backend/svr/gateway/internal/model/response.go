package model

import "oj-server/proto/pb"

type RankList struct {
	Phone     int64  `json:"phone"`
	NickName  string `json:"nickName"`
	PassCount int64  `json:"passCount"`
}

type SubmitRecordData struct {
	Uid         int64          `json:"uid,omitempty"`
	ProblemId   int64          `json:"problem_id,omitempty"`
	ProblemName string         `json:"problem_name,omitempty"`
	Status      string         `json:"status,omitempty"`
	Results     []*pb.PBResult `json:"results,omitempty"`
	Code        string         `json:"code,omitempty"`
	Lang        string         `json:"lang,omitempty"`

	ID       int64 `json:"id,omitempty"`
	CreateAt int64 `json:"create_at,omitempty"`
}
