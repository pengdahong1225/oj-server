package router

import (
	"github.com/gin-gonic/gin"
	"oj-server/svr/gateway/internal/handler"
	"oj-server/svr/gateway/internal/middlewares"
	"time"
)

func initUserRouter(rg *gin.RouterGroup) {
	// 用户
	user := rg.Group("/user")
	{
		user.POST("/login", handler.HandleUserLogin)
		user.POST("/login_sms", handler.HandleUserLoginBySms)
		user.POST("/refresh_token", handler.HandleReFreshAccessToken)
		user.POST("/register", handler.HandleUserRegister)
		user.POST("/reset_password", handler.HandleUserResetPassword)
		user.GET("/profile", middlewares.AuthLogin(), handler.HandleGetUserProfile)
	}

	// 验证码
	captcha := rg.Group("/captcha")
	captcha.Use(middlewares.RateLimitMiddleware(2*time.Second, 50))
	{
		captcha.GET("/image", handler.HandleGetImageCode)
		captcha.POST("/sms", handler.HandleGetSmsCode)
	}
}
