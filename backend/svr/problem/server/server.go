package server

import (
	"context"
	"fmt"
	"net"
	"oj-server/common/consts"
	"oj-server/common/proto/pb"
	"oj-server/pkg/go_pool"
	"oj-server/pkg/logger"
	"oj-server/pkg/registry"
	"oj-server/svr/problem/internal/client"
	"oj-server/svr/problem/internal/config"
	"oj-server/svr/problem/internal/service"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	reg        *registry.ConsulRegistry
	deregister func() error
	clients    *client.Clients

	listener       net.Listener
	problemService *service.ProblemService
	recordService  *service.RecordService
	commentService *service.CommentService
	noticeService  *service.NoticeService
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() error {
	configPath := fmt.Sprintf("%s/%s", consts.ConfigPath, "problem_config.yaml")
	err := config.MyConfig.ReadConfigFile(configPath)
	if err != nil {
		return err
	}
	err = logger.Init(consts.LogPath, config.MyConfig.Service.Name, logrus.DebugLevel)
	if err != nil {
		return err
	}

	dsn := fmt.Sprintf("%s:%d", config.MyConfig.Registry.Host, config.MyConfig.Registry.Port)
	s.reg, err = registry.NewConsulRegistry(dsn)
	if err != nil {
		logrus.Errorf("初始化注册中心失败: %v", err)
		return err
	}

	s.problemService, err = service.NewProblemService(s.clients)
	if err != nil {
		return err
	}
	s.recordService, err = service.NewRecordService(s.clients)
	if err != nil {
		return err
	}
	s.commentService, err = service.NewCommentService(s.clients)
	if err != nil {
		return err
	}
	s.noticeService, err = service.NewNoticeService(s.clients)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Run() error {
	serviceCfg := config.MyConfig.Service
	logrus.Infof("--------------- node_type:%v, node_id:%v, host:%v, port:%v, scheme:%v ---------------",
		serviceCfg.Name, serviceCfg.NodeId, serviceCfg.Address, serviceCfg.Port, serviceCfg.Scheme)

	deregister, err := s.reg.Register(context.Background(), registry.ServiceInfo{
		Name:    serviceCfg.Name,
		NodeId:  serviceCfg.NodeId,
		Address: serviceCfg.Address,
		Port:    serviceCfg.Port,
		Scheme:  serviceCfg.Scheme,
	})
	if err != nil {
		logrus.Errorf("注册服务失败: %v", err)
		return err
	}
	s.deregister = deregister

	s.clients, err = client.NewClients(s.reg)
	if err != nil {
		return fmt.Errorf("初始化客户端失败: %v", err)
	}

	// 建立排行榜
	go_pool.Submit(func() {
		s.recordService.SyncLeaderboardByScheduled()
	})
	// 消费judge-result队列
	go_pool.Submit(func() {
		s.recordService.ConsumeJudgeResult()
	})

	// 监听
	netAddr := fmt.Sprintf("%s:%d", config.MyConfig.Service.Address, config.MyConfig.Service.Port)
	listener, err := net.Listen("tcp", netAddr)
	if err != nil {
		logrus.Errorf("listen err: %v", err)
		return err
	}
	s.listener = listener

	// 启动GRPC服务
	grpcServer := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	pb.RegisterProblemServiceServer(grpcServer, s.problemService)
	pb.RegisterRecordServiceServer(grpcServer, s.recordService)
	pb.RegisterCommentServiceServer(grpcServer, s.commentService)
	pb.RegisterNoticeServiceServer(grpcServer, s.noticeService)
	return grpcServer.Serve(listener)
}

func (s *Server) Stop() {
	if s.deregister != nil {
		_ = s.deregister()
	}
	logrus.Warnf("======================= problem stop =======================")
}
