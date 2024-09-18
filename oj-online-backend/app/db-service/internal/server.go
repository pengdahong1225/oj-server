package internal

import (
	"fmt"
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/internal/logic"
	"github.com/pengdahong1225/Oj-Online-Server/common/goroutinePool"
	"github.com/pengdahong1225/Oj-Online-Server/common/registry"
	"github.com/pengdahong1225/Oj-Online-Server/common/settings"
	"github.com/pengdahong1225/Oj-Online-Server/common/utils"
	"github.com/pengdahong1225/Oj-Online-Server/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"sync"
	"time"
)

type Server struct {
}

func (receiver Server) Start() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorln(err)
		}
	}()

	wg := new(sync.WaitGroup)

	// DB服务
	wg.Add(1)
	err := goroutinePool.Instance().Submit(func() {
		defer wg.Done()
		StartRPCServer()
	})
	if err != nil {
		panic(err)
	}

	// 排行榜
	daemonServer := Daemon{}
	wg.Add(1)
	err = goroutinePool.Instance().Submit(func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				daemonServer.loopRank()
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

	wg.Wait()
}

func StartRPCServer() {
	system, err := settings.Instance().GetSystemConf("db-service")
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
	register, _ := registry.NewRegistry(settings.Instance().RegistryConfig)
	if err := register.RegisterService(system.Name, ip.String(), system.Port); err != nil {
		panic(err)
	}

	// 注册并启动db服务
	dbSrv := logic.DBServiceServer{}
	pb.RegisterDBServiceServer(grpcServer, &dbSrv)
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
