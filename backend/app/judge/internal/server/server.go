package server

import (
	"github.com/pengdahong1225/oj-server/backend/app/common/serverBase"
	"github.com/pengdahong1225/oj-server/backend/app/judge/internal/respository/cache"
	"github.com/pengdahong1225/oj-server/backend/app/judge/internal/service"
	"github.com/pengdahong1225/oj-server/backend/consts"
	"github.com/pengdahong1225/oj-server/backend/module/goroutinePool"
	"github.com/pengdahong1225/oj-server/backend/module/mq"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	ServerBase.Server
	judgeSrv *service.JudgeService
}

func (receiver *Server) Init() error {
	err := receiver.Initialize()
	if err != nil {
		return err
	}
	receiver.judgeSrv = service.NewJudgeService()
	err = cache.Init()
	if err != nil {
		return err
	}
	return nil
}

func (receiver *Server) Start() {
	go receiver.startJudgeConsume()

	err := receiver.Register()
	if err != nil {
		panic(err)
	}
}

func (receiver *Server) startJudgeConsume() {
	consumer := mq.NewConsumer(
		consts.RabbitMqExchangeKind,
		consts.RabbitMqExchangeName,
		consts.RabbitMqJudgeQueue,
		consts.RabbitMqJudgeKey,
		"", // 消费者标签，用于区别不同的消费者
	)
	deliveries := consumer.Consume()
	if deliveries == nil {
		logrus.Errorln("消费失败")
		return
	}
	defer consumer.Close()

	for d := range deliveries {
		if receiver.syncDo(d.Body) {
			d.Ack(false)
		} else {
			d.Reject(false)
		}
	}
}

// 解析，校验，提交任务给评测机
func (receiver *Server) syncDo(data []byte) bool {
	submitForm := &pb.SubmitForm{}
	err := proto.Unmarshal(data, submitForm)
	if err != nil {
		logrus.Errorln("解析err：", err.Error())
		return false
	}
	// 异步处理
	goroutinePool.Instance().Submit(func() {
		service.Handle(submitForm)
	})

	return true
}
