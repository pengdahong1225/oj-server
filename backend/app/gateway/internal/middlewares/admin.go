package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oj-server/module/auth"
)

func Admin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*auth.UserClaims)
		if currentUser.Authority == 0 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    1,
				"message": "无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
