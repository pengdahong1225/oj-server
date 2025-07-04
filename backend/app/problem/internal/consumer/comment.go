package consumer

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"oj-server/consts"
	"oj-server/module/mq"
	"oj-server/module/registry"
	"oj-server/proto/pb"
)

func ConsumeComment() {
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
	go writeComment(c)

	return true
}
func writeComment(c *pb.Comment) {
	conn, err := registry.NewDBConnection()
	if err != nil {
		logrus.Errorf("db-service连接失败, err=%v", err.Error())
		return
	}
	defer conn.Close()

	client := pb.NewCommentServiceClient(conn)
	_, err = client.SaveComment(context.Background(), &pb.SaveCommentRequest{
		Data: c,
	})
	if err != nil {
		logrus.Errorf("db-service保存评论失败, obj=%v,err=%v", c.ObjId, err.Error())
		return
	}
}
