package model

import "oj-server/pkg/proto/pb"

type Record struct {
	Uid         int64  `json:"uid"`
	UserName    string `json:"user_name"`
	ProblemID   int64  `json:"problem_id"`
	ProblemName string `json:"problem_name"`
	Status      string `json:"status"`
	Code        string `json:"code"`
	Result      []byte `json:"result"`
	Lang        string `json:"lang"`
}

func (r *Record) FromPbRecord(pbRecord *pb.SubmitRecord) {
	r.Uid = pbRecord.Uid
	r.UserName = pbRecord.UserName
	r.ProblemID = pbRecord.ProblemId
	r.ProblemName = pbRecord.ProblemName
	r.Status = pbRecord.Status
	r.Code = pbRecord.Code
	r.Result = pbRecord.Result
	r.Lang = pbRecord.Lang
}

// 查询用户历史提交记录参数
type QueryUserRecordListParams struct {
	Page     int32 `form:"page" binding:"required"`
	PageSize int32 `form:"page_size" binding:"required"`
}
type QueryUserRecordListResult struct {
	Total int64    `json:"total"`
	List  []Record `json:"list"`
}
