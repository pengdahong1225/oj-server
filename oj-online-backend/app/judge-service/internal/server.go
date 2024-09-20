package internal

import (
	"fmt"
	"github.com/pengdahong1225/Oj-Online-Server/app/judge-service/internal/judge"
	"github.com/pengdahong1225/Oj-Online-Server/app/judge-service/internal/svc/mq"
	"github.com/pengdahong1225/Oj-Online-Server/common/goroutinePool"
	"github.com/pengdahong1225/Oj-Online-Server/consts"
	"github.com/pengdahong1225/Oj-Online-Server/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"net/http"
	"sync"
)

type Server struct{}

func (receiver *Server) Loop(ip string, port int) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorln(err)
		}
	}()

	wg := new(sync.WaitGroup)
	// 健康检查线程
	wg.Add(1)
	goroutinePool.Instance().Submit(func() {
		defer wg.Done()
		dsn := fmt.Sprintf("%s:%d", ip, port)
		http.HandleFunc("/health", func(res http.ResponseWriter, req *http.Request) {
			if req.Method == "GET" {
				res.Write([]byte("ok"))
			}
		})
		http.ListenAndServe(dsn, nil)
	})

	// 消费者线程
	wg.Add(1)
	goroutinePool.Instance().Submit(func() {
		defer wg.Done()
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
	})
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
