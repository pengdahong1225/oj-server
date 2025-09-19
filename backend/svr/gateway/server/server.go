package server

import (
	"fmt"
	"oj-server/module/configs"
	"oj-server/module/registry"
	"oj-server/svr/gateway/internal/repo"
	"oj-server/svr/gateway/internal/router"

	"github.com/sirupsen/logrus"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() error {
	err := repo.Init()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Run() {
	// 服务注册
	err := registry.RegisterService()
	if err != nil {
		logrus.Fatalf("注册服务失败: %v", err)
	}

	// 启动
	server_cfg := configs.ServerConf
	engine := router.Router()
	dsn := fmt.Sprintf(":%d", server_cfg.Port)
	err = engine.Run(dsn)
	if err != nil {
		logrus.Errorf("%s", err)
		_ = registry.DeregisterService()
	}
}

func (s *Server) Stop() {
	_ = registry.DeregisterService()
	logrus.Errorf("======================= gateway stop =======================")
}
