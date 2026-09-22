package server

import (
	"context"
	"fmt"
	"oj-server/common/consts"
	"oj-server/pkg/logger"
	"oj-server/pkg/registry"
	"oj-server/svr/gateway/internal/api"
	"oj-server/svr/gateway/internal/client"
	"oj-server/svr/gateway/internal/config"

	"github.com/sirupsen/logrus"
)

type Server struct {
	reg        *registry.ConsulRegistry
	deregister func() error

	clients *client.Clients
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() error {
	configPath := fmt.Sprintf("%s/%s", consts.ConfigPath, "gateway_config.yaml")
	err := config.MyConfig.ReadConfigFile(configPath)
	if err != nil {
		return err
	}
	err = logger.Init(consts.LogPath, config.MyConfig.Service.Name, logrus.DebugLevel)
	if err != nil {
		return err
	}

	dsn := fmt.Sprintf("%s:%d", config.MyConfig.RegistryCfg.Host, config.MyConfig.RegistryCfg.Port)
	s.reg, err = registry.NewConsulRegistry(dsn)
	if err != nil {
		logrus.Errorf("初始化注册中心失败: %v", err)
		return err
	}

	return nil
}

func (s *Server) Run() error {
	serviceCfg := config.MyConfig.Service
	logrus.Infof("--------------- node_type:%v, node_id:%v, host:%v, port:%v, scheme:%v ---------------",
		serviceCfg.Name, serviceCfg.NodeId, serviceCfg.Address, serviceCfg.Port, serviceCfg.Scheme)

	deregister, err := s.reg.Register(context.Background(), registry.ServiceInfo{
		Name:    serviceCfg.Name,
		NodeId:  serviceCfg.NodeId,
		Address: serviceCfg.Address,
		Port:    serviceCfg.Port,
		Scheme:  serviceCfg.Scheme,
	})
	if err != nil {
		logrus.Errorf("注册服务失败: %v", err)
		return err
	}
	s.deregister = deregister

	s.clients, err = client.NewClients(s.reg)
	if err != nil {
		return fmt.Errorf("初始化客户端失败: %v", err)
	}

	engine := api.Router(s.clients)
	addr := fmt.Sprintf(":%d", serviceCfg.Port)
	return engine.Run(addr)
}

func (s *Server) Stop() {
	if s.deregister != nil {
		_ = s.deregister()
	}
	logrus.Warnf("======================= gateway stop =======================")
}
