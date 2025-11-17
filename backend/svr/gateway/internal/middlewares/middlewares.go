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
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// 可以加白名单校验
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers",
			"Content-Type, Authorization, X-Requested-With, X-CSRF-Token, Token, Accept, Origin")
		c.Header("Access-Control-Allow-Methods",
			"GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Expose-Headers",
			"Content-Length, Content-Type")
		c.Header("Access-Control-Max-Age", "86400") // 缓存 preflight

		// OPTIONS 请求直接返回，不继续执行 handler
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
