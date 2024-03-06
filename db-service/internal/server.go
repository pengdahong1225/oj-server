package internal

import (
	"db-service/global"
	"db-service/internal/handler"
	pb "db-service/proto"
	"db-service/utils"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

type Server struct {
}

func (receiver Server) Start() {
	// 获取ip地址
	ip, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	// 监听端口
	netAddr := fmt.Sprintf("%s:%d", ip.String(), global.ConfigInstance.System_.Port)
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
	registry, _ := global.NewRegistry()
	if err := registry.RegisterService(global.ConfigInstance.System_.Name, ip.String(), global.ConfigInstance.System_.Port); err != nil {
		panic(err)
	}

	// 注册并启动db服务
	dbSrv := handler.DBServiceServer{}
	pb.RegisterDBServiceServer(grpcServer, &dbSrv)
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
