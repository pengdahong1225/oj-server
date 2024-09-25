package internal

import (
	"fmt"
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/internal/daemon"
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/internal/rpc/comment"
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/internal/rpc/problem"
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/internal/rpc/record"
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/internal/rpc/user"
	"github.com/pengdahong1225/Oj-Online-Server/module/goroutinePool"
	"github.com/pengdahong1225/Oj-Online-Server/module/registry"
	"github.com/pengdahong1225/Oj-Online-Server/module/settings"
	"github.com/pengdahong1225/Oj-Online-Server/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"sync"
	"time"
)

type Server struct {
	Name string
	IP   string
	Port int
}

func (receiver *Server) Start() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorln(err)
		}
	}()

	wg := new(sync.WaitGroup)

	// DB rpc服务
	wg.Add(1)
	err := goroutinePool.Instance().Submit(func() {
		defer wg.Done()
		receiver.rpcServerStart()
	})
	if err != nil {
		panic(err)
	}

	// 排行榜
	daemonServer := daemon.Daemon{}
	wg.Add(1)
	err = goroutinePool.Instance().Submit(func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Hour * 24)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				daemonServer.LoopRank()
			}
		}
	})
	if err != nil {
		panic(err)
	}

	// mq消费者
	wg.Add(1)
	err = goroutinePool.Instance().Submit(func() {
		defer wg.Done()
		daemonServer.CommentSaveConsumer()
	})

	// 注册服务节点
	register, _ := registry.NewRegistry(settings.Instance().RegistryConfig)
	if err := register.RegisterServiceWithGrpc(receiver.Name, receiver.IP, receiver.Port); err != nil {
		panic(err)
	}

	wg.Wait()
}

func (receiver *Server) rpcServerStart() {
	// 监听端口
	netAddr := fmt.Sprintf("%s:%d", receiver.IP, receiver.Port)
	listener, err := net.Listen("tcp", netAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()

	// 启动rpc服务
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	userSrv := user.UserServer{}
	pb.RegisterUserServiceServer(grpcServer, &userSrv)
	problemSrv := problem.ProblemServer{}
	pb.RegisterProblemServiceServer(grpcServer, &problemSrv)
	recordSrv := record.RecordServer{}
	pb.RegisterRecordServiceServer(grpcServer, &recordSrv)
	commentSrv := comment.CommentServer{}
	pb.RegisterCommentServiceServer(grpcServer, &commentSrv)
	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
