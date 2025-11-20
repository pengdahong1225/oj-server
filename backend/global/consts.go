package global

import (
	"fmt"
	"time"
)

const (
	GatewayService = "gateway"
	UserService    = "user-service"
	ProblemService = "problem-service"
	JudgeService   = "judge-service"
)

const (
	ConfigPath                = "./config"
	LogPath                   = "./log"
	RefreshTokenTimeOut int64 = 60 * 60 * 24 * 7 // 7天
	AccessTokenTimeOut  int64 = 60 * 60 * 24     // 24小时
	Issuer                    = "Messi"
	ProblemConfigPath         = "/data/oj/problem_configs"
	GrpcTimeout               = time.Second * 3
	GrpcStreamTimeout         = time.Second * 30
)

// ============================ mq相关 ============================
const (
	RabbitMqExchangeKind = "direct"
	RabbitMqExchangeName = "amq.direct"

	RabbitMqJudgeSubmitKey   = "judge-submit"
	RabbitMqJudgeSubmitQueue = "judge-submit-queue"

	RabbitMqJudgeResultKey   = "judge-result"
	RabbitMqJudgeResultQueue = "judge-result-queue"

	RabbitMqCommentKey   = "comment"
	RabbitMqCommentQueue = "comment-task-queue"
)

// ============================ redis相关 ============================
const (
	// oj
	OjPrefix = "oj:"

	// tag list
	TagListKey = OjPrefix + "tag_list_set"

	// 短信验证码
	CaptchaPrefix  = OjPrefix + "captcha"
	CaptchaExpired = 60 * time.Second

	// 判题任务
	TaskStatePrefix   = OjPrefix + "task_state"
	TaskStateExpired  = 60 * 2 * time.Second     // 任务状态持续2min
	TaskResultPrefix  = OjPrefix + "task_result" // 判题结果前缀
	TaskResultExpired = 60 * time.Second         // 判题结果持续1min

	// 用户锁
	UserLockPrefix = OjPrefix + "user_lock" // 用户锁前缀
	UserLockTTL    = 60 * time.Second       // 用户锁TTL

	// 排行榜
	LeaderboardPrefix        = OjPrefix + "leaderboard:"
	AcTotalLeaderboardKey    = LeaderboardPrefix + "ac_total"
	LeaderboardLastUpdateKey = LeaderboardPrefix + "last_update"

	MonthLeaderboardTTL = time.Hour * 24 * 36 // 月榜key到期时间 36天
	DailyLeaderboardTTL = time.Hour * 24 * 2  // 日榜key到期时间 2天

	LeaderboardUserInfoKey = OjPrefix + "hot_user_info" // 排行榜用户信息
)

func GetMonthLeaderboardKey() string {
	return fmt.Sprintf("%s:%s", AcTotalLeaderboardKey, time.Now().Format("2006_01"))
}
func GetDailyLeaderboardKey() string {
	return fmt.Sprintf("%s:%s", AcTotalLeaderboardKey, time.Now().Format("2006_01_02"))
}
