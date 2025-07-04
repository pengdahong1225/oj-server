package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"oj-server/app/gateway/internal/handler"
	"oj-server/app/gateway/internal/middlewares"
	"os"
	"time"
)

func Router() *gin.Engine {
	path := fmt.Sprintf("%s/web.log", "./log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Errorf("web日志文件打开失败：%s", err.Error())
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout, file)
	gin.SetMode(os.Getenv("GIN_MODE"))

	r := gin.Default()
	r.Use(middlewares.Cors()) // 跨域处理

	// 初始化路由
	healthCheckRouters(r)
	problemRouters(r)

	return r
}

// healthCheckRouters 健康检查路由
func healthCheckRouters(engine *gin.Engine) {
	engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    "0",
			"message": "health",
		})
	})
}

// problemRouters
// 题目服务相关路由
// api/
func problemRouters(engine *gin.Engine) {
	// 排行榜
	engine.GET("/ranking_list", middlewares.AuthLogin(), handler.HandleGetRankList)

	// 验证码
	captchaRouter := engine.Group("/captcha")
	captchaRouter.Use(middlewares.RateLimitMiddleware(2*time.Second, 50))
	{
		captchaRouter.GET("/image", handler.HandleGetImageCode)
		captchaRouter.POST("/sms", handler.HandleGetSmsCode)
	}

	// 用户相关
	userRouter := engine.Group("/user")
	{
		userRouter.POST("/login", handler.HandleUserLogin)
		userRouter.POST("/login_sms", handler.HandleUserLoginBySms)
		userRouter.GET("/refresh_token", handler.HandleReFreshAccessToken)
		userRouter.POST("/register", handler.HandleUserRegister)
		userRouter.POST("/reset_password", handler.HandleUserResetPassword)
		userRouter.GET("/profile", middlewares.AuthLogin(), handler.HandleGetUserProfile)
		userRouter.GET("/record_list", middlewares.AuthLogin(), handler.HandleGetUserRecordList) // 历史提交记录列表
		userRouter.GET("/record", middlewares.AuthLogin(), handler.HandleGetUserRecord)
		userRouter.GET("/solved_list", middlewares.AuthLogin(), handler.HandleGetUserSolvedList)
	}

	// 题目相关
	problemRouter := engine.Group("/problem")
	{
		problemRouter.GET("/tag_list", handler.HandleGetTagList)
		problemRouter.GET("/list", handler.HandleGetProblemList)
		problemRouter.GET("/detail", handler.HandleGetProblemDetail)
		problemRouter.POST("/submit", middlewares.AuthLogin(), handler.HandleSubmitProblem)
		problemRouter.GET("/result", middlewares.AuthLogin(), handler.HandleGetSubmitResult) // 本次提交的结果
		problemRouter.POST("/add", middlewares.AuthLogin(), middlewares.Admin(), handler.HandleCreateProblem)
		problemRouter.POST("/update", middlewares.AuthLogin(), middlewares.Admin(), handler.HandleUpdateProblem)
		problemRouter.DELETE("", middlewares.AuthLogin(), middlewares.Admin(), handler.HandleDeleteProblem)
	}

	// 评论
	commentRouter := engine.Group("/comment")
	commentRouter.Use(middlewares.AuthLogin())
	{
		commentRouter.GET("/root_list", handler.HandleGetRootCommentList)
		commentRouter.GET("/child_list", handler.HandleGetChildCommentList)
		commentRouter.POST("/add", handler.HandleCreateComment)
		commentRouter.DELETE("", handler.HandleDeleteComment)
		commentRouter.POST("/like", handler.HandleLikeComment)
	}
}
