package service

import (
	"github.com/sirupsen/logrus"
	"oj-server/global"
	"oj-server/module/mq"
	"oj-server/module/proto/pb"
	"oj-server/svr/problem/internal/biz"
	"oj-server/svr/problem/internal/data"
)

// 题目服务
type ProblemService struct {
	pb.UnimplementedProblemServiceServer
	pc *biz.ProblemUseCase
	rc *biz.RecordUseCase

	problem_producer *mq.Producer // 判题任务生产者
	comment_consumer *mq.Consumer // 评论任务消费者
}

func NewProblemService() *ProblemService {
	var err error
	s := &ProblemService{}

	pr, err := data.NewProblemRepo()
	if err != nil {
		logrus.Fatalf("NewProblemService failed, err:%s", err.Error())
	}
	rr, err := data.NewRecordRepo()
	if err != nil {
		logrus.Fatalf("NewProblemService failed, err:%s", err.Error())
	}

	s.pc = biz.NewProblemUseCase(pr) // 注入实现
	s.rc = biz.NewRecordUseCase(rr)  // 注入实现

	s.problem_producer = mq.NewProducer(
		global.RabbitMqExchangeKind,
		global.RabbitMqExchangeName,
		global.RabbitMqJudgeQueue,
		global.RabbitMqJudgeKey,
	)
	s.comment_consumer = mq.NewConsumer(
		global.RabbitMqExchangeKind,
		global.RabbitMqExchangeName,
		global.RabbitMqCommentQueue,
		global.RabbitMqCommentKey,
		"", // 消费者标签，用于区别不同的消费者
	)

	return s
}
