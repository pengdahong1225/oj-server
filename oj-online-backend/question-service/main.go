package main

import (
	"fmt"
	"question-service/logger"
	"question-service/routers"
	"question-service/services/mq"
	"question-service/services/redis"
	"question-service/services/registry"
	"question-service/settings"
	"question-service/utils"
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
	// 初始化第三发服务
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		panic(err)
	}
	if err := mq.Init(settings.Conf.MqConfig); err != nil {
		panic(err)
	}
}

func Registry() {
	ip, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	registry, err := registry.NewRegistry(settings.Conf.RegistryConfig)
	if err != nil {
		panic(err)
	}
	if err = registry.RegisterService(settings.Conf.SystemConfig.Name, ip.String(), settings.Conf.SystemConfig.Port); err != nil {
		panic(err)
	}
}

func ServerLoop() {
	engine := routers.Router()
	dsn := fmt.Sprintf(":%d", settings.Conf.SystemConfig.Port)
	_ = engine.Run(dsn)
}

func main() {
	AppInit()
	Registry()
	ServerLoop()
}
