package service

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"oj-server/global"
	"oj-server/module/mq"
	"oj-server/proto/pb"
)

var (
	commentProducer *mq.Producer
)

func init() {
	commentProducer = mq.NewProducer(global.RabbitMqExchangeKind, global.RabbitMqExchangeName, global.RabbitMqCommentQueue, global.RabbitMqCommentKey)
}

func SendComment2MQ(comment *pb.Comment) error {
	msg, err := proto.Marshal(comment)
	if err != nil {
		logrus.Errorf("marshal comment failed: %v", err)
		return err
	}

	return commentProducer.Publish(msg)
}
