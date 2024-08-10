package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/middlewares"
	"os"
)

func Router() *gin.Engine {
	r := gin.Default()
	gin.SetMode(os.Getenv("GIN_MODE"))
	r.Use(middlewares.Cors()) // 跨域处理
	// 初始化路由
	HealthCheckRouters(r)
	QuestionRouters(r)
	CmsRouters(r)
	return r
}
