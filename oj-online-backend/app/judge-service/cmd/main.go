package main

import (
	"github.com/pengdahong1225/Oj-Online-Server/app/judge-service/internal"
	"github.com/pengdahong1225/Oj-Online-Server/module/logger"
	"github.com/pengdahong1225/Oj-Online-Server/module/registry"
	"github.com/pengdahong1225/Oj-Online-Server/module/settings"
	"github.com/pengdahong1225/Oj-Online-Server/module/utils"
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
	if err := logger.InitLog("judge-service", settings.Instance().LogConfig.Path, settings.Instance().LogConfig.Level); err != nil {
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
	system, err := settings.Instance().GetSystemConf("judge-service")
	if err != nil {
		panic(err)
	}
	register, err := registry.NewRegistry(settings.Instance().RegistryConfig)
	if err != nil {
		panic(err)
	}
	if err = register.RegisterServiceWithHttp(system.Name, ip.String(), system.Port); err != nil {
		panic(err)
	}

	// start
	server := &internal.Server{}
	server.Loop(ip.String(), system.Port)
}
