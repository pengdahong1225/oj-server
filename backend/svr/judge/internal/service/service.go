package service

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"oj-server/global"
	"oj-server/pkg/gPool"
	"oj-server/pkg/mq"
	"oj-server/pkg/proto/pb"
)

type JudgeService struct {
	problem_consumer *mq.Consumer // 判题任务消费者
	result_producer  *mq.Producer // 判题结果生产者
}

func NewJudgeService() *JudgeService {
	return &JudgeService{
		problem_consumer: mq.NewConsumer(
			global.RabbitMqExchangeKind,
			global.RabbitMqExchangeName,
			global.RabbitMqJudgeSubmitQueue,
			global.RabbitMqJudgeSubmitKey,
			"", // 消费者标签，用于区别不同的消费者
		),
		result_producer: mq.NewProducer(
			global.RabbitMqExchangeKind,
			global.RabbitMqExchangeName,
			global.RabbitMqJudgeResultQueue,
			global.RabbitMqJudgeResultKey,
		),
	}
}

func (s *JudgeService) ConsumeJudgeTask() {
	deliveries := s.problem_consumer.Consume()
	if deliveries == nil {
		logrus.Errorf("获取deliveries失败")
		return
	}
	defer s.problem_consumer.Close()

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
			_ = gPool.Instance().Submit(func() {
				s.Handle(task)
			})
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
