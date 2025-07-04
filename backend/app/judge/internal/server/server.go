package server

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"oj-server/app/common/serverBase"
	"oj-server/app/judge/internal/respository/cache"
	"oj-server/app/judge/internal/service"
	"oj-server/consts"
	"oj-server/module/goroutinePool"
	"oj-server/module/mq"
	"oj-server/proto/pb"
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
	err = goroutinePool.Instance().Submit(func() {
		service.HandleJudge(submitForm)
	})
	if err != nil {
		logrus.Errorln("异步处理err：", err.Error())
		return false
	}

	return true
}
