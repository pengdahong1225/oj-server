package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"net/http"
	"oj-server/global"
	"oj-server/module/configManager"
	"oj-server/module/gPool"
	"oj-server/module/mq"
	"oj-server/module/registry"
	"oj-server/proto/pb"
	"oj-server/src/judge/internal/service"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() error {
	err := service.Init()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Run() {
	go s.startJudgeConsume()

	// 服务注册
	err := registry.RegisterService()
	if err != nil {
		logrus.Fatalf("注册服务失败: %v", err)
	}

	cfg := configManager.ServerConf
	dsn := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	http.HandleFunc("/health", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			_, _ = res.Write([]byte("ok"))
		}
	})
	err = http.ListenAndServe(dsn, nil)
	if err != nil {
		logrus.Errorf("%s", err)
		_ = registry.DeregisterService()
	}
}

func (s *Server) Stop() {
	_ = registry.DeregisterService()
	logrus.Errorf("======================= judge stop =======================")
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
