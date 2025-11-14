package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oj-server/svr/gateway/internal/middlewares"
)

func Router() *gin.Engine {
	gin.SetMode(gin.DebugMode)

	engine := gin.New()
	engine.Use(middlewares.Logger(), middlewares.Recovery())
	engine.Use(middlewares.Cors()) // 跨域处理

	// 健康检查
	engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "health",
		})
	})

	v1 := engine.Group("/api/v1")
	// 用户服务路由
	initUserRouter(v1)

	// 题目服务路由
	initProblemRouter(v1)

	return engine
}
