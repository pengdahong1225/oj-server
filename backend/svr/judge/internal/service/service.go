package service

import (
	"oj-server/common/consts"
	"oj-server/common/proto/pb"
	"oj-server/pkg/mq"
	"oj-server/svr/judge/internal/client"
	"oj-server/svr/judge/internal/config"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type JudgeService struct {
	clients        *client.Clients
	jugeConsumer   *mq.Consumer // 判题任务消费者
	resultProducer *mq.Producer // 判题结果生产者
}

func NewJudgeService(clients *client.Clients) (*JudgeService, error) {
	cfg := config.MyConfig.MQ
	amqpClient := mq.NewClient(
		&mq.Options{
			Host:     cfg.Host,
			Port:     cfg.Port,
			User:     cfg.User,
			PassWord: cfg.PassWord,
			VHost:    cfg.VHost,
		},
	)

	return &JudgeService{
		jugeConsumer: &mq.Consumer{
			AmqpClient: amqpClient, // 注入client
			ExKind:     consts.RabbitMqExchangeKind,
			ExName:     consts.RabbitMqExchangeName,
			RoutingKey: consts.RabbitMqJudgeSubmitKey,
			QueName:    consts.RabbitMqJudgeSubmitQueue,
			CTag:       "", // 消费者标签，用于区别不同的消费者
		},
		resultProducer: &mq.Producer{
			AmqpClient: amqpClient, // 注入client
			ExKind:     consts.RabbitMqExchangeKind,
			ExName:     consts.RabbitMqExchangeName,
			QueName:    consts.RabbitMqJudgeResultQueue,
			RoutingKey: consts.RabbitMqJudgeResultKey,
		},
		clients: clients,
	}, nil
}

func (s *JudgeService) ConsumeJudgeTask() {
	deliveries := s.jugeConsumer.Consume()
	if deliveries == nil {
		logrus.Errorf("获取deliveries失败")
		return
	}
	defer s.jugeConsumer.Close()

	for d := range deliveries {
		// 处理任务
		result := func(data []byte) bool {
			task := new(pb.JudgeSubmission)
			err := proto.Unmarshal(data, task)
			if err != nil {
				logrus.Errorln("解析judge task err：", err.Error())
				return false
			}
			// 异步处理
			// _ = go_pool.Submit(func() {
			// 	s.Handle(task)
			// })
			return true
		}(d.Body)

		// 确认
		if result {
			_ = d.Ack(false)
		} else {
			_ = d.Reject(false)
		}
	}
}
