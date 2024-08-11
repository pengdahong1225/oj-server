package internal

import (
	"github.com/pengdahong1225/Oj-Online-Server/app/judge-service/internal/judge"
	"github.com/pengdahong1225/Oj-Online-Server/app/judge-service/services/mq"
	"github.com/pengdahong1225/Oj-Online-Server/common/goroutinePool"
	"github.com/pengdahong1225/Oj-Online-Server/consts"
	"github.com/pengdahong1225/Oj-Online-Server/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type Server struct{}

func (receiver *Server) Loop() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorln(err)
		}
	}()
	consumer := mq.NewConsumer(
		consts.RabbitMqExchangeKind,
		consts.RabbitMqExchangeName,
		consts.RabbitMqJudgeQueue,
		consts.RabbitMqJudgeKey,
		"", // 消费者标签，用于区别不同的消费者
	)
	deliveries := consumer.Consume()
	for d := range deliveries {
		if syncHandle(d.Body) {
			d.Ack(false)
		} else {
			d.Reject(false)
		}
	}

	select {}
}

// 解析，校验，提交任务给评测机
func syncHandle(data []byte) bool {
	submitForm := &pb.SubmitForm{}
	err := proto.Unmarshal(data, submitForm)
	if err != nil {
		logrus.Errorln("解析err：", err.Error())
		return false
	}
	// 处理
	goroutinePool.Instance().Submit(func() {
		judge.Handle(submitForm)
	})

	return true
}
