package consumer

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"oj-server/consts"
	"oj-server/module/gPool"
	"oj-server/module/mq"
	"oj-server/proto/pb"
)

func StartCommentConsume() {
	consumer := mq.NewConsumer(
		consts.RabbitMqExchangeKind,
		consts.RabbitMqExchangeName,
		consts.RabbitMqCommentQueue,
		consts.RabbitMqCommentKey,
		"", // 消费者标签，用于区别不同的消费者
	)
	deliveries := consumer.Consume()
	if deliveries == nil {
		logrus.Errorln("消费失败")
		return
	}
	defer consumer.Close()

	for d := range deliveries {
		if syncHandle(d.Body) {
			d.Ack(false)
		} else {
			d.Reject(false) // 拒绝并且丢失
		}
	}
}

func syncHandle(data []byte) bool {
	c := &pb.Comment{}
	err := proto.Unmarshal(data, c)
	if err != nil {
		logrus.Errorln("解析err：", err.Error())
		return false
	}
	// 处理
	_ = gPool.Instance().Submit(func() {
		writeComment(c)
	})

	return true
}
func writeComment(c *pb.Comment) {

}
