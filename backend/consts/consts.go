package consts

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
)

const (
	GatewayService = "gateway-service"
	UserService    = "user-service"
	ProblemService = "problem-service"
	JudgeService   = "judge-service"
)
