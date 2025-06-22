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
	TokenTimeOut int64 = 60 * 60 * 24 * 7
	Issuer             = "Messi"
)

const (
	GatewayService = "gateway-service"
	UserService    = "user-service"
	ProblemService = "problem-service"
	JudgeService   = "judge-service"
)
