package routers

import (
	"github.com/gin-gonic/gin"
	"os"
	"question-service/middlewares"
)

func Router() *gin.Engine {
	r := gin.Default()
	gin.SetMode(os.Getenv("GIN_MODE"))
	r.Use(middlewares.Cors())
	HealthCheckRouters(r)
	QuestionRouters(r)
	CmsRouters(r)
	return r
}
