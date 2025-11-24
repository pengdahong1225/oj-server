package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"oj-server/global"
	"oj-server/pkg/logger"
	"oj-server/pkg/proto/pb"
	"oj-server/pkg/registry"
	"oj-server/svr/problem/internal/configs"
	"oj-server/svr/problem/internal/service"
)

type Server struct {
	listener       net.Listener
	problemService *service.ProblemService
	recordService  *service.RecordService
	commentService *service.CommentService
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() error {
	// 初始化日志
	server_cfg := configs.ServerConf
	err := logger.Init(global.LogPath, server_cfg.Name, logrus.DebugLevel)
	if err != nil {
		return err
	}

	// 初始化注册中心
	registry_cfg := configs.AppConf.RegistryCfg
	dsn := fmt.Sprintf("%s:%d", registry_cfg.Host, registry_cfg.Port)
	registrar, err := registry.NewRegistrar(dsn)
	if err != nil {
		logrus.Errorf("初始化注册中心失败: %v", err)
		return err
	}
	registry.MyRegistrar = registrar

	s.problemService = service.NewProblemService()
	s.recordService = service.NewRecordService()
	s.commentService = service.NewCommentService()

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

	// 建立排行榜
	go s.recordService.SyncLeaderboardByScheduled()

	// 消费judge-result队列
	go s.recordService.ConsumeJudgeResult()

	// 监听
	cfg := configs.ServerConf
	netAddr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	listener, err := net.Listen("tcp", netAddr)
	if err != nil {
		logrus.Fatalf("listen err: %v", err)
		return
	}
	s.listener = listener

	// 启动GRPC服务
	grpcServer := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	pb.RegisterProblemServiceServer(grpcServer, s.problemService)
	pb.RegisterRecordServiceServer(grpcServer, s.recordService)
	pb.RegisterCommentServiceServer(grpcServer, s.commentService)
	_ = grpcServer.Serve(listener)
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
	logrus.Warnf("======================= problem stop =======================")
}
