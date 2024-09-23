package main

import (
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/internal"
	"github.com/pengdahong1225/Oj-Online-Server/module/logger"
	"github.com/pengdahong1225/Oj-Online-Server/module/settings"
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

func ServerLoop() {
	server := internal.Server{}
	server.Start()
}

func main() {
	AppInit()
	ServerLoop()
}
