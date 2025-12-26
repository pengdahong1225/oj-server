package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oj-server/svr/gateway/internal/model"
)

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
