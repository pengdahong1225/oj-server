package internal

import (
	"fmt"
	ServerBase "github.com/pengdahong1225/oj-server/backend/app/common/serverBase"
	"github.com/pengdahong1225/oj-server/backend/app/judge-service/internal/judge"
	"github.com/pengdahong1225/oj-server/backend/consts"
	"github.com/pengdahong1225/oj-server/backend/module/goroutinePool"
	"github.com/pengdahong1225/oj-server/backend/module/mq"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"net/http"
)

type Server struct {
	ServerBase.Server
}

func (receiver *Server) Start() {
	goroutinePool.Instance().Submit(func() {
		dsn := fmt.Sprintf("%s:%d", receiver.Host, receiver.Port)
		http.HandleFunc("/health", func(res http.ResponseWriter, req *http.Request) {
			if req.Method == "GET" {
				res.Write([]byte("ok"))
			}
		})
		http.ListenAndServe(dsn, nil)
		logrus.Errorf("健康检查线程退出")
	})

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
		if syncDo(d.Body) {
			d.Ack(false)
		} else {
			d.Reject(false)
		}
	}
}

// 解析，校验，提交任务给评测机
func syncDo(data []byte) bool {
	submitForm := &pb.SubmitForm{}
	err := proto.Unmarshal(data, submitForm)
	if err != nil {
		logrus.Errorln("解析err：", err.Error())
		return false
	}
	// 异步处理
	goroutinePool.Instance().Submit(func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorln("panic:", r)
			}
		}()
		judge.Handle(submitForm)
	})

	return true
}
