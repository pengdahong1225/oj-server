package internal

import (
	"fmt"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/routers"
	"github.com/pengdahong1225/Oj-Online-Server/module/registry"
	"github.com/pengdahong1225/Oj-Online-Server/module/settings"
)

type Server struct {
	Name string
	IP   string
	Port int
}

func (receiver *Server) Register() error {
	register, err := registry.NewRegistry(settings.Instance().RegistryConfig)
	if err != nil {
		return err
	}
	if err = register.RegisterServiceWithHttp(receiver.Name, receiver.IP, receiver.Port); err != nil {
		return err
	}
	return nil
}
func (receiver *Server) Start() {
	engine := routers.Router()
	dsn := fmt.Sprintf(":%d", receiver.Port)
	_ = engine.Run(dsn)
}
