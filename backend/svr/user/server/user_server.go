package server

import (
	"fmt"
	"net"
	"oj-server/module/configManager"
	"oj-server/module/registry"
	"oj-server/proto/pb"
	"oj-server/src/user/internal/service"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
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
	// 服务注册
	err := registry.RegisterService()
	if err != nil {
		logrus.Fatalf("注册服务失败: %v", err)
	}

	// 监听
	cfg := configManager.ServerConf
	netAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	listener, err := net.Listen("tcp", netAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	s.listener = listener

	// 启动
	grpcServer := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	pb.RegisterUserServiceServer(grpcServer, s.userSrv)
	err = grpcServer.Serve(listener)
	if err != nil {
		logrus.Errorf("%s", err)
		_ = registry.DeregisterService()
	}
}

func (s *Server) Stop() {
	_ = s.listener.Close()
	logrus.Errorf("======================= user stop =======================")
}
