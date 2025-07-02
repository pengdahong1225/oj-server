package define

import "github.com/pengdahong1225/oj-server/backend/proto/pb"

// Response 统一返回格式
type Response struct {
	Code    int    `json:"code"` // 业务状态码
	Message string `json:"message"`
	Data    any    `json:"data"`
}

const (
	Success             = 0
	Failed              = 1
	AccessTokenExpired  = 2 // 访问令牌过期
	RefreshTokenExpired = 3 // 刷新令牌过期
	Unauthorized        = 4 // 未登录
	TokenInvalid        = 5 // 令牌无效
)

type RankList struct {
	Phone     int64  `json:"phone"`
	NickName  string `json:"nickName"`
	PassCount int64  `json:"passCount"`
}

type LoginRspData struct {
	Rsp         *pb.UserLoginResponse `json:"data"`
	AccessToken string                `json:"access_token"`
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
