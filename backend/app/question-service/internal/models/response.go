package models

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
