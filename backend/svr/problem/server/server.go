package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"oj-server/module/configs"
	"oj-server/module/registry"
	"oj-server/proto/pb"
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
	s.problemService = service.NewProblemService()
	s.recordService = service.NewRecordService()
	s.commentService = service.NewCommentService()

	return nil
}

func (s *Server) Run() {
	// 评论任务消费
	go s.commentService.ConsumeComment()
	// 异步维护排行榜
	go s.recordService.UpdateLeaderboardByScheduled()

	// 服务注册
	err := registry.RegisterService()
	if err != nil {
		logrus.Fatalf("注册服务失败: %v", err)
	}
	// 监听
	cfg := configs.ServerConf
	netAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	listener, err := net.Listen("tcp", netAddr)
	if err != nil {
		logrus.Fatalf("监听失败: %v", err)
	}
	defer listener.Close()
	s.listener = listener

	grpcServer := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	pb.RegisterProblemServiceServer(grpcServer, s.problemService)
	pb.RegisterRecordServiceServer(grpcServer, s.recordService)
	pb.RegisterCommentServiceServer(grpcServer, s.commentService)
	if err = grpcServer.Serve(listener); err != nil {
		logrus.Errorf("%s", err)
		_ = registry.DeregisterService()
	}
}

func (s *Server) Stop() {
	_ = s.listener.Close()
	logrus.Errorf("======================= problem stop =======================")
}
