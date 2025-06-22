package define

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
