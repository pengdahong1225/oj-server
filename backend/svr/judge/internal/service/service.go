package service

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"oj-server/global"
	"oj-server/pkg/gPool"
	"oj-server/pkg/mq"
	"oj-server/pkg/proto/pb"
	"oj-server/svr/judge/internal/biz"
	"oj-server/svr/judge/internal/data"
)

type JudgeService struct {
	uc       *biz.JudgeUseCase
	consumer *mq.Consumer
}

func NewJudgeService() *JudgeService {
	repo, err := data.NewRepo()
	if err != nil {
		logrus.Fatalf("new repo failed, err:%v", err)
	}
	uc := biz.NewJudgeUseCase(repo)

	return &JudgeService{
		uc: uc,
		consumer: mq.NewConsumer(
			global.RabbitMqExchangeKind,
			global.RabbitMqExchangeName,
			global.RabbitMqJudgeQueue,
			global.RabbitMqJudgeKey,
			""), // 消费者标签，用于区别不同的消费者
	}
}

func (s *JudgeService) ConsumeJudgeTask() {
	deliveries := s.consumer.Consume()
	if deliveries == nil {
		logrus.Errorf("获取deliveries失败")
		return
	}
	defer s.consumer.Close()

	for d := range deliveries {
		// 处理任务
		result := func(data []byte) bool {
			task := &pb.SubmitForm{}
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
