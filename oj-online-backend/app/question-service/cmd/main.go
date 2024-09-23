package main

import (
	"fmt"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/routers"
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
	if err := logger.InitLog("question-service", settings.Instance().LogConfig.Path, settings.Instance().LogConfig.Level); err != nil {
		panic(err)
	}
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
	system, err := settings.Instance().GetSystemConf("question-service")
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
	// loop
	ServerLoop(system.Port)
}
