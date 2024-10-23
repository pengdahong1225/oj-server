package main

import (
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/svc/redis"
	"github.com/pengdahong1225/oj-server/backend/module/logger"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"github.com/pengdahong1225/oj-server/backend/module/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := logger.InitLog("db-service", settings.Instance().LogConfig.Path, settings.Instance().LogConfig.Level); err != nil {
		panic(err)
	}
	redis.Init(settings.Instance().RedisConfig.Host, settings.Instance().RedisConfig.Port)

	system, err := settings.Instance().GetSystemConf("db-service")
	if err != nil {
		logrus.Errorln(err.Error())
		panic(err)
	}
	// 获取ip地址
	ip, err := utils.GetOutboundIP()
	if err != nil {
		logrus.Errorln(err.Error())
		panic(err)
	}
	server := internal.Server{
		Name: system.Name,
		IP:   ip.String(),
		Port: system.Port,
	}
	server.Start()
}
