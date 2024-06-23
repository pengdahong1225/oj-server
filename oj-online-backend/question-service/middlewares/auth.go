package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// jwt鉴权头部信息 x-token 登录时返回token信息
		token := ctx.Request.Header.Get("token")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "未登录",
			})
			ctx.Abort()
			return
		}
		j := NewJWT()
		// 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, TokenExpired) {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"code":    http.StatusUnauthorized,
					"message": "授权已过期",
				})
				ctx.Abort()
				return
			} else {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"code":    http.StatusUnauthorized,
					"message": "token验证失败",
				})
				ctx.Abort()
				return
			}
		}
		// token通过
		ctx.Set("claims", claims)
		ctx.Next()
	}
}
