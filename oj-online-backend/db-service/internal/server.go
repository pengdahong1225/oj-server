package internal

import (
	"db-service/internal/daemon"
	"db-service/internal/handler"
	"db-service/internal/proto"
	"db-service/services/ants"
	"db-service/services/registry"
	"db-service/settings"
	"db-service/utils"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"sync"
)

type Server struct {
}

func (receiver Server) Start() {
	// 后台-排行榜服务
	wg := new(sync.WaitGroup)
	wg.Add(2)
	err := ants.AntsPoolInstance.Submit(func() {
		defer wg.Done()
		daemon.StartDaemon()
	})
	if err != nil {
		panic(err)
	}
	// DB服务
	err = ants.AntsPoolInstance.Submit(func() {
		defer wg.Done()
		StartRPCServer()
	})
	if err != nil {
		panic(err)
	}
	wg.Wait()
}

func StartRPCServer() {
	system, err := settings.GetSystemConf("db-service")
	if err != nil {
		panic(err)
	}
	// 获取ip地址
	ip, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	// 监听端口
	netAddr := fmt.Sprintf("%s:%d", ip.String(), system.Port)
	listener, err := net.Listen("tcp", netAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()
	// 健康检查
	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthcheck)
	// 注册
	register, _ := registry.NewRegistry(settings.Conf.RegistryConfig)
	if err := register.RegisterService(system.Name, ip.String(), system.Port); err != nil {
		panic(err)
	}

	// 注册并启动db服务
	dbSrv := handler.DBServiceServer{}
	pb.RegisterDBServiceServer(grpcServer, &dbSrv)
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
