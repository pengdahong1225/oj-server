package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"oj-server/global"
	"oj-server/pkg/logger"
	"oj-server/pkg/registry"
	"oj-server/svr/gateway/internal/configs"
	"oj-server/svr/gateway/internal/router"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() error {
	// 初始化日志
	server_cfg := configs.ServerConf
	err := logger.Init(global.LogPath, server_cfg.Name, logrus.DebugLevel)
	if err != nil {
		return err
	}

	// 初始化注册中心
	registry_cfg := configs.AppConf.RegistryCfg
	dsn := fmt.Sprintf("redis://%s:%d", registry_cfg.Host, registry_cfg.Port)
	registrar, err := registry.NewRegistrar(dsn)
	if err != nil {
		logrus.Errorf("初始化注册中心失败: %v", err)
		return err
	}
	registry.MyRegistrar = registrar

	return nil
}

func (s *Server) Run() {
	server_cfg := configs.ServerConf
	logrus.Infof("--------------- node_type:%v, node_id:%v, host:%v, port:%v, scheme:%v ---------------",
		server_cfg.Name, server_cfg.NodeId, server_cfg.Address, server_cfg.Port, server_cfg.Scheme)

	// 注册服务
	err := registry.MyRegistrar.RegisterService(&registry.ServiceInfo{
		Name:    server_cfg.Name,
		NodeId:  server_cfg.NodeId,
		Address: server_cfg.Address,
		Port:    server_cfg.Port,
		Scheme:  server_cfg.Scheme,
	})
	if err != nil {
		logrus.Fatalf("注册服务失败: %v", err)
		return
	}

	// 启动HTTP服务
	engine := router.Router()
	if err = engine.Run(fmt.Sprintf("%s:%d", server_cfg.Address, server_cfg.Port)); err != nil {
		logrus.Fatalf("启动HTTP服务失败: %v", err)
	}
}

func (s *Server) Stop() {
	server_cfg := configs.ServerConf
	_ = registry.MyRegistrar.UnRegister(&registry.ServiceInfo{
		Name:    server_cfg.Name,
		NodeId:  server_cfg.NodeId,
		Address: server_cfg.Address,
		Port:    server_cfg.Port,
		Scheme:  server_cfg.Scheme,
	})
	logrus.Warnf("======================= gateway stop =======================")
}
