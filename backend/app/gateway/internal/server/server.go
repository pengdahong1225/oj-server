package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"oj-server/app/common/serverBase"
	"oj-server/app/gateway/internal/respository/cache"
	"oj-server/app/gateway/internal/router"
	"sync"
)

// Server
// 服务器
type Server struct {
	ServerBase.Server
	wg sync.WaitGroup
}

func (s *Server) Init() error {
	err := s.Initialize()
	if err != nil {
		return err
	}
	err = cache.Init()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Start() {
	// 启动
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		engine := router.Router()
		dsn := fmt.Sprintf(":%d", s.Port)
		err := engine.Run(dsn)
		if err != nil {
			logrus.Fatalf("启动服务失败: %v", err)
		}
	}()

	// 服务注册
	err := s.Register()
	if err != nil {
		logrus.Fatalf("注册服务失败: %v", err)
	}
	defer s.UnRegister()

	s.wg.Wait()
}
