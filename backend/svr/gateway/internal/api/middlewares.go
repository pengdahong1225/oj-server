package api

import (
	"errors"
	"net/http"
	"oj-server/pkg/jwt_utils"
	"oj-server/svr/gateway/internal/config"
	"oj-server/svr/gateway/internal/model"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"github.com/sirupsen/logrus"
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

// 令牌桶其实和漏桶的原理类似，令牌桶按固定的速率往桶里放入令牌，并且只要能从桶里取出令牌就能通过，令牌桶支持突发流量的快速处理。
// 对于从桶里取不到令牌的场景，我们可以选择等待也可以直接拒绝并返回。
// 对于令牌桶的Go语言实现，大家可以参照github.com/juju/ratelimit库。这个库支持多种令牌桶模式，并且使用起来也比较简单。

// 对于该限流中间件的注册位置，我们可以按照不同的限流策略将其注册到不同的位置，例如：
// 如果要对全站限流就可以注册成全局的中间件。
// 如果是某一组路由需要限流，那么就只需将该限流中间件注册到对应的路由组即可。

// RateLimitMiddleware 创建指定填充速率和容量大小的令牌桶
func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		// 如果取不到令牌就中断本次请求返回 rate limit...
		if bucket.TakeAvailable(1) == 0 {
			c.String(http.StatusOK, "服务繁忙")
			c.Abort()
			return
		}
		// 取到令牌就放行
		c.Next()
	}
}
func AuthLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未登录",
			})
			ctx.Abort()
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		// 解析token包含的信息
		jwt_builder := jwt_utils.JWTBuilder{
			SigningKey: []byte(config.MyConfig.JwtCfg.SigningKey), // 密钥
		}
		claims := new(model.UserClaims)

		err := jwt_builder.ParseToken(token, claims)
		if err != nil {
			if errors.Is(err, jwt_utils.TokenExpired) {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "token已过期",
				})
				ctx.Abort()
				return
			} else {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "token错误",
				})
				ctx.Abort()
				return
			}
		}

		// 剩余时间
		remain := claims.ExpiresAt - time.Now().Unix()
		logrus.Debugf("剩余：%v, [%v天]", remain, remain/86400)

		// token通过
		ctx.Set("claims", claims)
		ctx.Set("uid", claims.Uid)
		ctx.Set("role", claims.Authority)
		ctx.Next()
	}
}

func Admin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value, exists := ctx.Get("claims")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无登录信息",
			})
			ctx.Abort()
			return
		}
		userClaims, ok := value.(*model.UserClaims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无登录信息",
			})
			ctx.Abort()
			return
		}
		if userClaims.Authority == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
