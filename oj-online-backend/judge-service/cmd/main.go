package main

import (
	"github.com/pengdahong1225/common-utils"
	"judge-service/internal"
	"judge-service/services/goroutinePool"
	"judge-service/services/judgeClient"
	"judge-service/services/logger"
	"judge-service/services/redis"
	"judge-service/services/registry"
	"judge-service/settings"
	"time"
)

func AppInit() {
	// 配置全局时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = loc
	// 初始化
	if err := settings.Init(); err != nil {
		panic(err)
	}
	if err := logger.Init(); err != nil {
		panic(err)
	}
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		panic(err)
	}
	if err := goroutinePool.Init(); err != nil {
		panic(err)
	}
	if err := judgeClient.Init(); err != nil {
		panic(err)
	}
}

func main() {
	// 初始化
	AppInit()
	// 注册服务节点
	ip, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	system, err := settings.GetSystemConf("judge-service")
	if err != nil {
		panic(err)
	}
	register, err := registry.NewRegistry(settings.Conf.RegistryConfig)
	if err != nil {
		panic(err)
	}
	if err = register.RegisterService(system.Name, ip.String(), system.Port); err != nil {
		panic(err)
	}

	// start
	server := &internal.Server{}
	server.Loop()
}
