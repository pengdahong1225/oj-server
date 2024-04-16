package main

import (
	"db-service/internal"
	"db-service/logger"
	ants2 "db-service/services/ants"
	"db-service/services/dao/mysql"
	"db-service/services/dao/redis"
	"db-service/services/registry"
	"db-service/settings"
	"db-service/utils"
	"github.com/panjf2000/ants/v2"
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
	if err := mysql.Init(settings.Conf.MysqlConfig); err != nil {
		panic(err)
	}
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
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
	server := internal.Server{}
	server.Start()
}

func main() {
	AppInit()
	Registry()
	ServerLoop()
}
