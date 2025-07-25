package consumer

import (
	"oj-server/global"
	"oj-server/module/gPool"
	"oj-server/module/mq"
	"oj-server/proto/pb"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

var (
	// 评论任务消费者
	comment_consumer *mq.Consumer
)

func init() {
	comment_consumer = mq.NewConsumer(
		global.RabbitMqExchangeKind,
		global.RabbitMqExchangeName,
		global.RabbitMqCommentQueue,
		global.RabbitMqCommentKey,
		"", // 消费者标签，用于区别不同的消费者
	)
}

func StartCommentConsume() {
	deliveries := comment_consumer.Consume()
	if deliveries == nil {
		logrus.Errorln("消费失败")
		return
	}
	defer comment_consumer.Close()

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
