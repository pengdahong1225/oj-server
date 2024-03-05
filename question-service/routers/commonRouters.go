package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 健康检查路由
func HealthCheckRouters(engine *gin.Engine) {
	engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "health",
		})
	})
}
