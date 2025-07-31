package global

import "time"

const (
	GatewayService = "gateway"
	UserService    = "user"
	ProblemService = "problem"
	JudgeService   = "judge"
)

const (
	ConfigPath                = "./config"
	LogPath                   = "./log"
	RefreshTokenTimeOut int64 = 60 * 60 * 24 * 7
	AccessTokenTimeOut  int64 = 60 * 15
	Issuer                    = "Messi"
	ProblemConfigPath         = "/data/problem_config"
)

// ============================ mq相关 ============================
const (
	RabbitMqExchangeKind = "direct"
	RabbitMqExchangeName = "amq.direct"

	RabbitMqJudgeKey     = "judge"
	RabbitMqJudgeQueue   = "judge-task-queue"
	RabbitMqCommentKey   = "comment"
	RabbitMqCommentQueue = "comment-task-queue"
)

// ============================ redis相关 ============================
const (
	// 短信验证码
	CaptchaPrefix  = "captcha"
	CaptchaExpired = 60 * time.Second

	// 判题任务
	TaskStatePrefix   = "task_state"
	TaskStateExpired  = 60 * 2 * time.Second // 任务状态持续2min
	TaskResultPrefix  = "task_result"        // 判题结果前缀
	TaskResultExpired = 60 * time.Second     // 判题结果持续1min

	// 用户锁
	UserLockPrefix = "user_lock"      // 用户锁前缀
	UserLockTTL    = 60 * time.Second // 用户锁TTL
)
