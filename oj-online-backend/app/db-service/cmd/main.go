package main

import (
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/internal"
	"github.com/pengdahong1225/Oj-Online-Server/common/logger"
	"github.com/pengdahong1225/Oj-Online-Server/common/registry"
	"github.com/pengdahong1225/Oj-Online-Server/common/settings"
	"github.com/pengdahong1225/Oj-Online-Server/common/utils"
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
	if err := logger.InitLog("db-service", settings.Instance().LogConfig.Path, settings.Instance().LogConfig.Level); err != nil {
		panic(err)
	}
}

func Registry() {
	ip, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	register, err := registry.NewRegistry(settings.Instance().RegistryConfig)
	if err != nil {
		panic(err)
	}
	system, err := settings.Instance().GetSystemConf("db-service")
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
