package main

import (
	"github.com/pengdahong1225/Oj-Online-Server/app/judge-service/internal"
	"github.com/pengdahong1225/Oj-Online-Server/module/logger"
	"github.com/pengdahong1225/Oj-Online-Server/module/settings"
	"github.com/pengdahong1225/Oj-Online-Server/module/utils"
	"time"
)

func main() {
	// 初始化
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = loc
	if err := logger.InitLog("judge-service", settings.Instance().LogConfig.Path, settings.Instance().LogConfig.Level); err != nil {
		panic(err)
	}

	system, err := settings.Instance().GetSystemConf("judge-service")
	if err != nil {
		panic(err)
	}
	ip, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	// start
	server := &internal.Server{
		Name: system.Name,
		IP:   ip.String(),
		Port: system.Port,
	}
	server.Start()
}
