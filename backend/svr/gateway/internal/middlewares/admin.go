package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oj-server/module/auth"
	"oj-server/proto/pb"
	"oj-server/svr/gateway/internal/api/define"
)

func Admin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*auth.UserClaims)
		if currentUser.Authority == 0 {
			resp := &define.Response{
				ErrCode: pb.Error_EN_Unauthorized,
				Message: "无权限",
			}
			ctx.JSON(http.StatusUnauthorized, resp)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
