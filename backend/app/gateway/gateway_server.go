package gateway

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"oj-server/app/gateway/internal/respository/cache"
	"oj-server/app/gateway/internal/router"
	"oj-server/module/configManager"
	"oj-server/module/registry"
	"sync"
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
	server_cfg := configManager.ServerConf
	wg := sync.WaitGroup{}
	// 启动
	wg.Add(1)
	go func() {
		defer wg.Done()
		engine := router.Router()
		dsn := fmt.Sprintf(":%d", server_cfg.Port)
		err := engine.Run(dsn)
		if err != nil {
			logrus.Fatalf("启动服务失败: %v", err)
		}
	}()

	// 服务注册
	err := registry.RegisterService()
	if err != nil {
		logrus.Fatalf("注册服务失败: %v", err)
	}
	defer func() {
		err = registry.DeregisterService()
		if err != nil {
			logrus.Errorf("注销服务失败: %v", err)
			return
		}
	}()

	wg.Wait()
}
