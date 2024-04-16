package main

import (
	"github.com/panjf2000/ants/v2"
	"judge-service/internal"
	"judge-service/logger"
	ants2 "judge-service/services/ants"
	"judge-service/services/mq"
	"judge-service/services/redis"
	"judge-service/services/registry"
	"judge-service/settings"
	"judge-service/utils"
	"time"
)

func AppInit() {
	// 配置全局时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = loc
	// 初始化配置
	if err := settings.Init(); err != nil {
		panic(err)
	}
	// 初始化日志
	if err := logger.Init(); err != nil {
		panic(err)
	}
	// 初始化服务组件
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		panic(err)
	}
	if err := mq.Init(settings.Conf.MqConfig); err != nil {
		panic(err)
	}
	ants2.AntsPoolInstance, _ = ants.NewPool(ants.DefaultAntsPoolSize, ants.WithPanicHandler(ants2.AntsPanicHandler))
}

func Registry() {
	ip, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	register, err := registry.NewRegistry(settings.Conf.RegistryConfig)
	if err != nil {
		panic(err)
	}
	if err = register.RegisterService(settings.Conf.SystemConfig.Name, ip.String(), settings.Conf.SystemConfig.Port); err != nil {
		panic(err)
	}
}

func ServerLoop() {
	srv := internal.Server{}
	srv.Start()
}

func main() {
	AppInit()
	Registry()
	ServerLoop()
}
