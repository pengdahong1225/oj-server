package internal

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pengdahong1225/oj-server/backend/app/common/serverBase"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"sync"
)

// Server
// 服务器
type Server struct {
	ServerBase.Server
	wg sync.WaitGroup
}

func (receiver *Server) Init() error {
	err := receiver.Initialize()
	if err != nil {
		return err
	}

	return nil
}

func (receiver *Server) Start() {
	mux := runtime.NewServeMux()

	mux.Handle

	receiver.registerGrpcServices(mux, "problem-service", pb.RegisterProblemServiceHandlerFromEndpoint)
	receiver.registerGrpcServices(mux, "judge-service", pb.RegisterJudgeServiceHandlerFromEndpoint)
	receiver.registerGrpcServices(mux, "user-service", pb.RegisterUserServiceHandlerFromEndpoint)

	receiver.wg.Add(1)
	go func() {
		defer receiver.wg.Done()
		err := http.ListenAndServe(":8080", mux)
		if err != nil {
			logrus.Fatalf("启动服务失败: %v", err)
		}
	}()

	// 服务注册
	err := receiver.Register()
	if err != nil {
		logrus.Fatalf("注册服务失败: %v", err)
	}
	defer receiver.UnRegister()

	receiver.wg.Wait()
}

func (receiver *Server) registerGrpcServices(mux *runtime.ServeMux, serviceName string, handlerFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error) {
	// 从服务发现获取地址
	addr := discoverService(serviceName)
	err := handlerFunc(
		context.Background(),
		mux,
		addr,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		},
	)
	if err != nil {
		logrus.Fatalf("注册服务 %s 失败: %v", serviceName, err)
	}
}
