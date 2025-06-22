package internal

import (
	"fmt"
	"github.com/pengdahong1225/oj-server/backend/app/common/serverBase"
	"github.com/pengdahong1225/oj-server/backend/app/problem-service/internal/biz/middlewares"
	"github.com/pengdahong1225/oj-server/backend/app/problem-service/internal/biz/service"
	"github.com/pengdahong1225/oj-server/backend/app/problem-service/internal/consumer"
	"github.com/pengdahong1225/oj-server/backend/app/problem-service/internal/repository/cache"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"sync"
)

type Server struct {
	ServerBase.Server
	problemSrv *service.ProblemService
	wg         sync.WaitGroup
}

func (receiver *Server) Init() error {
	err := receiver.Initialize()
	if err != nil {
		return err
	}

	receiver.problemSrv = service.NewProblemService()
	err = cache.Init()
	if err != nil {
		return err
	}

	return nil
}

func (receiver *Server) Start() {
	// 拦截器
	chain := grpc.ChainUnaryInterceptor(
		middlewares.RecoveryInterceptor,
		middlewares.LoggingInterceptor,
	)
	grpcServer := grpc.NewServer(chain)

	// 健康检查
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	// 题目服务
	problemSrv := service.ProblemService{}
	pb.RegisterProblemServiceServer(grpcServer, &problemSrv)

	// 启动
	receiver.wg.Add(1)
	go func() {
		defer receiver.wg.Done()
		netAddr := fmt.Sprintf("%s:%d", receiver.Host, receiver.Port)
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

	// 启动消费者
	go consumer.ConsumeComment()

	err := receiver.Register()
	if err != nil {
		panic(err)
	}
	defer receiver.UnRegister()

	receiver.wg.Wait()
}
