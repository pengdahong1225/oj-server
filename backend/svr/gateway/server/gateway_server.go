package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"oj-server/global"
	"oj-server/module/configs"
	"oj-server/module/registry"
	"oj-server/proto/pb"
	"oj-server/svr/gateway/internal/middlewares"
)

type Server struct {
	engine *gin.Engine       // http 引擎
	gwMux  *runtime.ServeMux // gRPC-Gateway多路复用器
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() error {
	gin.SetMode(gin.ReleaseMode)
	s.engine = gin.Default()

	s.gwMux = runtime.NewServeMux()

	return nil
}

func (s *Server) Run() {
	err := registry.RegisterService()
	if err != nil {
		logrus.Fatalf("注册服务失败: %v", err)
	}
	defer registry.DeregisterService()

	// 服务注册表
	services := []struct {
		name    string
		handler func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
	}{
		{global.UserService, pb.RegisterUserServiceHandler},
		{global.ProblemService, pb.RegisterProblemServiceHandler},
		{global.ProblemService, pb.RegisterRecordServiceHandler},
		{global.ProblemService, pb.RegisterCommentServiceHandler},
	}

	// 动态连接并注册所有服务
	var conns []*grpc.ClientConn
	for _, svc := range services {
		conn, err := registry.GetGrpcConnection(svc.name)
		if err != nil {
			logrus.Errorf("failed to connect to %s: %v", svc.name, err)
			continue
		}
		conns = append(conns, conn)

		if err = svc.handler(context.Background(), s.gwMux, conn); err != nil {
			logrus.Errorf("failed to register %s handler: %v", svc.name, err)
		} else {
			logrus.Infof("successfully registered %s", svc.name)
		}
	}

	// 将所有其他请求路由到 gRPC-Gateway
	s.engine.Any("/*", gin.WrapH(s.gwMux))
	s.SetupMiddlewares()

	addr := fmt.Sprintf("%s:%d", configs.ServerConf.Host, configs.ServerConf.Port)
	logrus.Infof("starting gateway on %s", addr)
	err = s.engine.Run(addr)
	if err != nil {
		logrus.Errorf("failed to serve: %v", err)
	}
}

func (s *Server) SetupMiddlewares() {
	s.engine.Use(
		middlewares.Cors(),
	)
}

func (s *Server) Stop() {
	_ = registry.DeregisterService()
	logrus.Errorf("======================= gateway stop =======================")
}
