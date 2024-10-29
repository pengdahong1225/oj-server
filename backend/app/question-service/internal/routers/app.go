package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/middlewares"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func Router() *gin.Engine {
	path := fmt.Sprintf("%s/web.log", settings.Instance().LogConfig.Path)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Errorf("web日志文件打开失败：%s", err.Error())
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout, file)
	gin.SetMode(os.Getenv("GIN_MODE"))

	r := gin.Default()
	r.Use(middlewares.Cors()) // 跨域处理

	// 初始化路由
	healthCheckRouters(r)
	questionRouters(r)

	return r
}
