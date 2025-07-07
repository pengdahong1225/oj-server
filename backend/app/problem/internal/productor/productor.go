package productor

import (
	"oj-server/consts"
	"oj-server/module/mq"
)

var (
	// 判题任务生产者
	problem_productor *mq.Producer
)

func init() {
	problem_productor = mq.NewProducer(
		consts.RabbitMqExchangeKind,
		consts.RabbitMqExchangeName,
		consts.RabbitMqJudgeQueue,
		consts.RabbitMqJudgeKey,
	)
}

func Publish(data []byte) bool {
	return problem_productor.Publish(data)
}
