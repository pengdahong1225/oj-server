package server

import (
	"fmt"
	"net"
	"oj-server/pkg/registry"
	"oj-server/svr/user/internal/service"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"oj-server/pkg/proto/pb"
	"oj-server/svr/user/internal/configs"
)

type Server struct {
	listener net.Listener
	userSrv  *service.UserService
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() error {
	s.userSrv = service.NewUserService()
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

	// 监听
	netAddr := fmt.Sprintf("%s:%d", server_cfg.Address, server_cfg.Port)
	listener, err := net.Listen("tcp", netAddr)
	if err != nil {
		logrus.Fatalf("listen err: %v", err)
		return
	}
	s.listener = listener

	// 启动GRPC服务
	grpcServer := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	pb.RegisterUserServiceServer(grpcServer, s.userSrv)
	if err = grpcServer.Serve(listener); err != nil {
		logrus.Fatalf("启动GRPC服务失败: %v", err)
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
	_ = s.listener.Close()
	logrus.Warnf("======================= user stop =======================")
}
