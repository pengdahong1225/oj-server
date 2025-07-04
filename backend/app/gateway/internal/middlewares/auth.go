package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"oj-server/module/auth"
	"oj-server/module/settings"
	"time"
)

func AuthLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// jwt鉴权头部信息 token
		token := ctx.Request.Header.Get("access-token")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    1,
				"message": "未登录",
			})
			ctx.Abort()
			return
		}
		signingKey := settings.Instance().SigningKey
		j := auth.JWTCreator{
			SigningKey: []byte(signingKey),
		}
		// 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, auth.TokenExpired) {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"code":    1,
					"message": "授权已过期",
				})
				ctx.Abort()
				return
			} else {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"code":    1,
					"message": "token验证失败",
				})
				ctx.Abort()
				return
			}
		}

		// 剩余时间
		remain := claims.ExpiresAt - time.Now().Unix()
		logrus.Debugf("剩余：%v, [%v天]\n", remain, remain/86400)

		// token通过
		ctx.Set("claims", claims)
		ctx.Set("uid", claims.Uid)
		ctx.Next()
	}
}
