package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"oj-server/svr/gateway/internal/configs"
	"strings"
	"time"
	"oj-server/pkg/jwt_utils"
	"oj-server/svr/gateway/internal/model"
)

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
			SigningKey: []byte(configs.AppConf.JwtCfg.SigningKey), // 密钥
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
