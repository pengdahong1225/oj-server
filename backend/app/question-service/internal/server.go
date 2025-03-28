package internal

import (
	"fmt"
	"github.com/pengdahong1225/oj-server/backend/app/common/serverBase"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/cache"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/consumer"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/routers"
)

type Server struct {
	ServerBase.Server
}

func (receiver *Server) Init() error {
	err := receiver.Initialize()
	if err != nil {
		return err
	}

	err = cache.Init()
	if err != nil {
		return err
	}

	return nil
}

func (receiver *Server) Start() {
	engine := routers.Router()
	dsn := fmt.Sprintf(":%d", receiver.Port)
	go func() {
		err := engine.Run(dsn)
		if err != nil {
			panic(err)
		}
	}()

	// 启动消费者
	go consumer.ConsumeComment()

	err := receiver.Register()
	if err != nil {
		panic(err)
	}

	select {}
}
