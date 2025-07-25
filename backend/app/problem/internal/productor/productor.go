package productor

import (
	"oj-server/global"
	"oj-server/module/mq"
)

var (
	// 判题任务生产者
	problem_productor *mq.Producer
)

func init() {
	problem_productor = mq.NewProducer(
		global.RabbitMqExchangeKind,
		global.RabbitMqExchangeName,
		global.RabbitMqJudgeQueue,
		global.RabbitMqJudgeKey,
	)
}

func Publish(data []byte) bool {
	return problem_productor.Publish(data)
}
