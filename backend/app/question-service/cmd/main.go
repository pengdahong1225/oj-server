package main

import (
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/svc/cache"
	"github.com/pengdahong1225/oj-server/backend/module/logger"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"github.com/pengdahong1225/oj-server/backend/module/utils"
)

func main() {
	if err := logger.InitLog("db-service", settings.Instance().LogConfig.Path, settings.Instance().LogConfig.Level); err != nil {
		panic(err)
	}
	cache.Init(settings.Instance().RedisConfig.Host, settings.Instance().RedisConfig.Port)

	system, err := settings.Instance().GetSystemConf("question-service")
	if err != nil {
		panic(err)
	}
	ip, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	// start
	server := internal.Server{
		Name: system.Name,
		IP:   ip.String(),
		Port: system.Port,
	}
	if err = server.Register(); err != nil {
		panic(err)
	}
	server.Start()
}
