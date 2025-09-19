package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"oj-server/module/configs"
	"oj-server/module/registry"
	"oj-server/svr/judge/internal/service"
)

type Server struct {
	judgeService *service.JudgeService
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() error {
	s.judgeService = service.NewJudgeService()
	if s.judgeService == nil {
		return fmt.Errorf("judge service init failed")
	}

	return nil
}

func (s *Server) Run() {
	// 启动判题任务消费
	go s.judgeService.ConsumeJudgeTask()

	// 服务注册
	err := registry.RegisterService()
	if err != nil {
		logrus.Fatalf("注册服务失败: %v", err)
	}
	cfg := configs.ServerConf
	dsn := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	http.HandleFunc("/health", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			_, _ = res.Write([]byte("ok"))
		}
	})
	if err = http.ListenAndServe(dsn, nil); err != nil {
		logrus.Errorf("%s", err)
		_ = registry.DeregisterService()
	}
}

func (s *Server) Stop() {
	_ = registry.DeregisterService()
	logrus.Errorf("======================= judge stop =======================")
}
