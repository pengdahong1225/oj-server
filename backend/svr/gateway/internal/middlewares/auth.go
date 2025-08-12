package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"oj-server/module/auth"
	"oj-server/module/configManager"
	"oj-server/proto/pb"
	"oj-server/src/gateway/internal/define"
	"time"
)

func AuthLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp := &define.Response{}

		// jwt鉴权头部信息 token
		token := ctx.Request.Header.Get("access-token")
		if token == "" {
			resp.ErrCode = pb.Error_EN_Unauthorized
			resp.Message = "未登录"
			ctx.JSON(http.StatusUnauthorized, resp)
			ctx.Abort()
			return
		}
		signingKey := configManager.AppConf.JwtCfg.SigningKey
		j := auth.JWTCreator{
			SigningKey: []byte(signingKey),
		}
		// 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, auth.TokenExpired) {
				resp.ErrCode = pb.Error_EN_AccessTokenExpired
				resp.Message = "授权已过期"
				ctx.JSON(http.StatusUnauthorized, resp)
				ctx.Abort()
				return
			} else {
				resp.ErrCode = pb.Error_EN_TokenInvalid
				resp.Message = "token验证失败"
				ctx.JSON(http.StatusUnauthorized, resp)
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
