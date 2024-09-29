package main

import (
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal"
	"github.com/pengdahong1225/oj-server/backend/module/logger"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"github.com/pengdahong1225/oj-server/backend/module/utils"
	"time"
)

func main() {
	// 初始化
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = loc
	if err := logger.InitLog("db-service", settings.Instance().LogConfig.Path, settings.Instance().LogConfig.Level); err != nil {
		panic(err)
	}

	system, err := settings.Instance().GetSystemConf("db-service")
	if err != nil {
		panic(err)
	}
	// 获取ip地址
	ip, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	server := internal.Server{
		Name: system.Name,
		IP:   ip.String(),
		Port: system.Port,
	}
	server.Start()
}
