package consts

import "time"

const (
	GatewayService = "gateway-service"
	UserService    = "user-service"
	ProblemService = "problem-service"
	JudgeService   = "judge-service"
)

const (
	RabbitMqExchangeKind = "direct"
	RabbitMqExchangeName = "amq.direct"

	RabbitMqJudgeKey     = "judge"
	RabbitMqJudgeQueue   = "judge-task-queue"
	RabbitMqCommentKey   = "comment"
	RabbitMqCommentQueue = "comment-task-queue"
)

const (
	RefreshTokenTimeOut int64 = 60 * 60 * 24 * 7
	AccessTokenTimeOut  int64 = 60 * 15
	Issuer                    = "Messi"
	ProblemConfigPath         = "/data/problem_config"
	TaskStateExpired          = 60 * 2 * time.Second // 任务状态持续2min
	JudgeResultExpired        = 60 * time.Second     // 判题结果持续1min
)
