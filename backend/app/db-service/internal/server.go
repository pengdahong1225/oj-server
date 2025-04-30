package internal

import (
	"fmt"
	ServerBase "github.com/pengdahong1225/oj-server/backend/app/common/serverBase"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/api/comment"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/api/notice"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/api/problem"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/api/record"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/api/user"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/cache"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/mysql"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

type Server struct {
	ServerBase.Server
}

func (receiver *Server) Init() error {
	err := receiver.Initialize()
	if err != nil {
		return err
	}
	err = mysql.Init()
	if err != nil {
		return err
	}
	err = cache.Init()
	if err != nil {
		return err
	}

	return nil
}

func (receiver *Server) Start() {
	netAddr := fmt.Sprintf("%s:%d", receiver.Host, receiver.Port)
	listener, err := net.Listen("tcp", netAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()

	// 健康检查
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	//  服务
	userSrv := user.UserServer{}
	pb.RegisterUserServiceServer(grpcServer, &userSrv)
	problemSrv := problem.ProblemServer{}
	pb.RegisterProblemServiceServer(grpcServer, &problemSrv)
	recordSrv := record.RecordServer{}
	pb.RegisterRecordServiceServer(grpcServer, &recordSrv)
	commentSrv := comment.CommentServer{}
	pb.RegisterCommentServiceServer(grpcServer, &commentSrv)
	noticeSrv := notice.NoticeServer{}
	pb.RegisterNoticeServiceServer(grpcServer, &noticeSrv)

	go func() {
		err = grpcServer.Serve(listener)
		if err != nil {
			panic(err)
		}
	}()

	err = receiver.Register()
	if err != nil {
		panic(err)
	}

	select {}
}
