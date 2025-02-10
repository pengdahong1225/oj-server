package internal

import (
	"fmt"
	"github.com/pengdahong1225/oj-server/backend/app/common/serverBase"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/routers"
)

type Server struct {
	ServerBase.Server
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

	err := receiver.Register()
	if err != nil {
		panic(err)
	}

	select {}
}
