package gateway

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"oj-server/app/gateway/internal/respository/cache"
	"oj-server/app/gateway/internal/router"
	"oj-server/module/configManager"
	"oj-server/module/registry"
)

// Server
// 服务器
type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() error {
	err := cache.Init()
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
	server_cfg := configManager.ServerConf
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
