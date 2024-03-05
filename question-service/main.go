package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"question-service/global"
	"question-service/middlewares"
	"question-service/routers"
	"question-service/utils"
	"time"
)

func main() {
	// 配置全局时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = loc

	// 注册
	ip, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	registry, err := global.NewRegistry()
	if err != nil {
		panic(err)
	}
	if err = registry.RegisterService(global.ConfigInstance.System_.Name, ip.String(), global.ConfigInstance.System_.Port); err != nil {
		panic(err)
	}

	// gin
	engine := gin.Default()
	engine.Use(middlewares.Cors())
	routers.HealthCheckRouters(engine)
	routers.QuestionRouters(engine)
	dsn := fmt.Sprintf(":%d", global.ConfigInstance.System_.Port)
	_ = engine.Run(dsn)
}
