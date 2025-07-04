package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"oj-server/app/common/serverBase"
	"oj-server/app/user/internal/respository/cache"
	"oj-server/app/user/internal/service"
	"oj-server/proto/pb"
	"sync"
)

type Server struct {
	ServerBase.Server
	userSrv *service.UserService
	wg      sync.WaitGroup
}

func (s *Server) Init() error {
	err := s.Initialize()
	if err != nil {
		return err
	}

	s.userSrv = service.NewUserService()
	err = cache.Init()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Start() {
	// tls认证
	creds, err := credentials.NewServerTLSFromFile("./config/keys/server.crt", "./config/keys/server.key")
	if err != nil {
		logrus.Fatalf("Failed to generate credentials %v", err)
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.Creds(creds))
	grpcServer := grpc.NewServer(opts...)

	// 健康检查
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	// 将业务服务注册到grpc中
	pb.RegisterUserServiceServer(grpcServer, s.userSrv)

	// 启动
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		netAddr := fmt.Sprintf("%s:%d", s.Host, s.Port)
		listener, err := net.Listen("tcp", netAddr)
		if err != nil {
			panic(err)
		}
		defer listener.Close()
		err = grpcServer.Serve(listener)
		if err != nil {
			logrus.Fatalf("启动服务失败: %v", err)
		}
	}()

	err = s.Register()
	if err != nil {
		panic(err)
	}
	defer s.UnRegister()

	s.wg.Wait()
}
