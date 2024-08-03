package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"question-service/internal/routers"
	ants2 "question-service/services/ants"
	"question-service/services/judgeService"
	"question-service/services/logger"
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
	// 初始化服务组件
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		panic(err)
	}
	ants2.AntsPoolInstance, _ = ants.NewPool(ants.DefaultAntsPoolSize, ants.WithPanicHandler(ants2.AntsPanicHandler))
	judgeService.Init()
}

func ServerLoop(port int) {
	engine := routers.Router()
	dsn := fmt.Sprintf(":%d", port)
	_ = engine.Run(dsn)
}

func main() {
	// 初始化
	AppInit()
	// 注册服务节点
	ip, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	system, err := settings.GetSystemConf("question-service")
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
	// loop
	ServerLoop(system.Port)
}
