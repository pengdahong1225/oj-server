package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// [web] | 200 |     3.63421ms |   192.168.2.116 | POST     "/topology/generate"
func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		// 执行请求
		ctx.Next()

		// 计算耗时
		latency := time.Since(start)
		statusCode := ctx.Writer.Status()
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		path := ctx.Request.URL.Path

		// 打印日志
		logrus.Infof("[web] | %d | %10v | %15s | %-8s \"%s\"",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}
func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("panic: %v", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "Internal Server Error",
				})
			}
		}()
		ctx.Next()
	}
}

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method

		ctx.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization")
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access_Control-Allow-Origin, Access_Control-Headers, Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
	}
}
