package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Admin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*UserClaims)
		if currentUser.Authority == 0 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Set("uid", currentUser.Uid)
		ctx.Next()
	}
}
