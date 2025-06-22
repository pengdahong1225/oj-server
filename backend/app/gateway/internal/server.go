package internal

import (
	"fmt"
	"github.com/pengdahong1225/oj-server/backend/app/common/serverBase"
	"github.com/pengdahong1225/oj-server/backend/app/gateway/internal/router"
	"github.com/sirupsen/logrus"
	"sync"
)

// Server
// 服务器
type Server struct {
	ServerBase.Server
	wg sync.WaitGroup
}

func (receiver *Server) Init() error {
	err := receiver.Initialize()
	if err != nil {
		return err
	}

	return nil
}

func (receiver *Server) Start() {
	// 启动
	receiver.wg.Add(1)
	go func() {
		defer receiver.wg.Done()
		engine := router.Router()
		dsn := fmt.Sprintf(":%d", receiver.Port)
		err := engine.Run(dsn)
		if err != nil {
			logrus.Fatalf("启动服务失败: %v", err)
		}
	}()

	// 服务注册
	err := receiver.Register()
	if err != nil {
		logrus.Fatalf("注册服务失败: %v", err)
	}
	defer receiver.UnRegister()

	receiver.wg.Wait()
}
