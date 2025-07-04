package productor

import (
	"oj-server/consts"
	"oj-server/module/mq"
)

var (
	productor *mq.Producer
)

func init() {
	productor = mq.NewProducer(
		consts.RabbitMqExchangeKind,
		consts.RabbitMqExchangeName,
		consts.RabbitMqJudgeQueue,
		consts.RabbitMqJudgeKey,
	)
}

func Publish(data []byte) bool {
	return productor.Publish(data)
}
