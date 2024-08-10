package main

import (
	"github.com/pengdahong1225/Oj-Online-Server/app/judge-service/internal"
	"github.com/pengdahong1225/Oj-Online-Server/app/judge-service/services/goroutinePool"
	"github.com/pengdahong1225/Oj-Online-Server/app/judge-service/services/judgeClient"
	"github.com/pengdahong1225/Oj-Online-Server/app/judge-service/services/redis"
	"github.com/pengdahong1225/Oj-Online-Server/app/judge-service/settings"
	"github.com/pengdahong1225/Oj-Online-Server/config"
	"github.com/pengdahong1225/Oj-Online-Server/pkg/logger"
	"github.com/pengdahong1225/Oj-Online-Server/pkg/registry"
	"github.com/pengdahong1225/Oj-Online-Server/pkg/utils"
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
	if err := logger.InitLog("judge-service", settings.Conf.LogConfig.Path, settings.Conf.LogConfig.Level); err != nil {
		panic(err)
	}
	if err := goroutinePool.Init(); err != nil {
		panic(err)
	}
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
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
	system, err := config.GetSystemConf(settings.Conf.SystemConfigs, "judge-service")
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
