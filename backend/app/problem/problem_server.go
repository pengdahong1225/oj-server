package problem

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"oj-server/app/problem/internal/consumer"
	"oj-server/app/problem/internal/repository/cache"
	"oj-server/app/problem/internal/service"
	"oj-server/module/configManager"
	"oj-server/module/registry"
	"oj-server/proto/pb"
	"sync"
)

type Server struct {
	problemSrv *service.ProblemService
}

func (s *Server) Init() error {
	s.problemSrv = service.NewProblemService()
	err := cache.Init()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Run() {
	var opts []grpc.ServerOption
	// tls认证
	creds, err := credentials.NewServerTLSFromFile("./config/keys/server.pem", "./config/keys/server.key")
	if err != nil {
		logrus.Fatalf("Failed to generate credentials %v", err)
	}
	opts = append(opts, grpc.Creds(creds))

	grpcServer := grpc.NewServer(opts...)

	// 健康检查
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	// 将业务服务注册到grpc中
	pb.RegisterProblemServiceServer(grpcServer, s.problemSrv)

	// 启动
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		cfg := configManager.ServerConf
		netAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
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

	// 启动评论消费者
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.StartCommentConsume()
	}()

	// 服务注册
	err = registry.RegisterService()
	if err != nil {
		logrus.Fatalf("注册服务失败: %v", err)
	}
	defer func() {
		err = registry.DeregisterService()
		if err != nil {
			logrus.Errorf("注销服务失败: %v", err)
			return
		}
	}()

	wg.Wait()
}
