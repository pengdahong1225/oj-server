package judge

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"oj-server/app/judge/internal/respository/cache"
	"oj-server/app/judge/internal/service"
	"oj-server/global"
	"oj-server/module/gPool"
	"oj-server/module/mq"
	"oj-server/module/registry"
	"oj-server/proto/pb"
	"sync"
)

type Server struct{}

func (s *Server) Init() error {
	err := service.Init()
	if err != nil {
		return err
	}
	err = cache.Init()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Run() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.startJudgeConsume()
	}()

	// 服务注册
	err := registry.RegisterService()
	if err != nil {
		logrus.Fatalf("注册服务失败: %v", err)
	}
	defer func() {
		err = registry.DeregisterService()
		if err != nil {
			logrus.Errorf("注销服务失败: %v", err)
			return
		}
	}()

	wg.Wait()
}

func (s *Server) startJudgeConsume() {
	consumer := mq.NewConsumer(
		global.RabbitMqExchangeKind,
		global.RabbitMqExchangeName,
		global.RabbitMqJudgeQueue,
		global.RabbitMqJudgeKey,
		"", // 消费者标签，用于区别不同的消费者
	)
	deliveries := consumer.Consume()
	if deliveries == nil {
		logrus.Errorln("消费失败")
		return
	}
	defer consumer.Close()

	for d := range deliveries {
		if s.syncDo(d.Body) {
			d.Ack(false)
		} else {
			d.Reject(false)
		}
	}
}

// 解析，校验，提交任务给评测机
func (s *Server) syncDo(data []byte) bool {
	submitForm := &pb.SubmitForm{}
	err := proto.Unmarshal(data, submitForm)
	if err != nil {
		logrus.Errorln("解析err：", err.Error())
		return false
	}
	// 异步处理
	err = gPool.Instance().Submit(func() {
		service.Handle(submitForm)
	})
	if err != nil {
		logrus.Errorln("异步处理err：", err.Error())
		return false
	}

	return true
}
