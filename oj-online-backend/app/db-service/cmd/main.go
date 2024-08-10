package main

import (
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/internal"
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/services/mysql"
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/services/redis"
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/settings"
	"github.com/pengdahong1225/Oj-Online-Server/app/judge-service/services/goroutinePool"
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
	if err := logger.InitLog("db-service", settings.Conf.LogConfig.Path, settings.Conf.LogConfig.Level); err != nil {
		panic(err)
	}
	if err := mysql.Init(settings.Conf.MysqlConfig); err != nil {
		panic(err)
	}
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		panic(err)
	}
	if err := goroutinePool.Init(); err != nil {
		panic(err)
	}
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
	system, err := settings.GetSystemConf("db-service")
	if err != nil {
		panic(err)
	}
	if err = register.RegisterService(system.Name, ip.String(), system.Port); err != nil {
		panic(err)
	}
}

func ServerLoop() {
	server := internal.Server{}
	server.Start()
}

func main() {
	AppInit()
	Registry()
	ServerLoop()
}
