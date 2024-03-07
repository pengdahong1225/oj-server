package routers

import (
	"github.com/gin-gonic/gin"
	"question-service/middlewares"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())
	HealthCheckRouters(r)
	QuestionRouters(r)
	CmsRouters(r)
	return r
}
