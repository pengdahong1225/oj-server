package internal

import (
	"fmt"
	"github.com/pengdahong1225/oj-server/backend/app/judge-service/internal/judge"
	"github.com/pengdahong1225/oj-server/backend/consts"
	"github.com/pengdahong1225/oj-server/backend/module/goroutinePool"
	"github.com/pengdahong1225/oj-server/backend/module/mq"
	"github.com/pengdahong1225/oj-server/backend/module/registry"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"net/http"
	"sync"
)

type Server struct {
	Name string
	IP   string
	Port int
	UUID string
}

func (receiver *Server) register() error {
	register, err := registry.NewRegistry(settings.Instance().RegistryConfig)
	if err != nil {
		return err
	}
	if err = register.RegisterServiceWithHttp(receiver.Name, receiver.IP, receiver.Port, receiver.UUID); err != nil {
		return err
	}
	return nil
}

func (receiver *Server) Start() {
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
		dsn := fmt.Sprintf("%s:%d", receiver.IP, receiver.Port)
		http.HandleFunc("/health", func(res http.ResponseWriter, req *http.Request) {
			if req.Method == "GET" {
				res.Write([]byte("ok"))
			}
		})
		http.ListenAndServe(dsn, nil)
		logrus.Infoln("健康检查线程退出")
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
		if deliveries == nil {
			logrus.Errorln("消费失败")
			return
		}
		defer consumer.Close()

		for d := range deliveries {
			if syncHandle(d.Body) {
				d.Ack(false)
			} else {
				d.Reject(false)
			}
		}
		logrus.Infoln("消费者线程退出")
	})

	if err := receiver.register(); err != nil {
		panic(err)
	}

	wg.Wait()
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
