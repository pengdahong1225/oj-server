package main

import (
	"fmt"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/routers"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/services/goroutinePool"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/services/mq"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/services/redis"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/setting"
	"github.com/pengdahong1225/Oj-Online-Server/pkg/logger"
	"github.com/pengdahong1225/Oj-Online-Server/pkg/registry"
	"github.com/pengdahong1225/Oj-Online-Server/pkg/settings"
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
	if err := logger.InitLog("question-service", setting.Instance().LogConfig.Path, setting.Instance().LogConfig.Level); err != nil {
		panic(err)
	}
	if err := redis.Init(setting.Instance().RedisConfig); err != nil {
		panic(err)
	}
	if err := goroutinePool.Init(); err != nil {
		panic(err)
	}
	if err := mq.Init(); err != nil {
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
	system, err := settings.GetSystemConf(setting.Instance().SystemConfigs, "question-service")
	if err != nil {
		panic(err)
	}
	register, err := registry.NewRegistry(setting.Instance().RegistryConfig)
	if err != nil {
		panic(err)
	}
	if err = register.RegisterService(system.Name, ip.String(), system.Port); err != nil {
		panic(err)
	}
	// loop
	ServerLoop(system.Port)
}
