package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"oj-server/module/configManager"
	"oj-server/module/registry"
	"oj-server/proto/pb"
	"oj-server/svr/problem/internal/service"
)

type Server struct {
	listener   net.Listener
	problemSrv *service.ProblemService
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() error {
	s.problemSrv = service.NewProblemService()

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
		logrus.Fatalf("监听失败: %v", err)
	}
	defer listener.Close()
	s.listener = listener

	// 启动
	go s.problemSrv.StartCommentConsume()

	grpcServer := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	pb.RegisterProblemServiceServer(grpcServer, s.problemSrv)
	err = grpcServer.Serve(listener)
	if err != nil {
		logrus.Errorf("%s", err)
		_ = registry.DeregisterService()
	}
}

func (s *Server) Stop() {
	_ = s.listener.Close()
	logrus.Errorf("======================= problem stop =======================")
}
