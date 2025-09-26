package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"oj-server/global"
	"oj-server/svr/gateway/internal/handler"
	"oj-server/svr/gateway/internal/middlewares"
	"os"
	"time"
)

func Router() *gin.Engine {
	timer := time.Now().Format("2006_01_02")
	path := fmt.Sprintf("%s/web.%s.log", global.LogPath, timer)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Errorf("web日志文件打开失败：%s", err.Error())
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout, file)
	gin.SetMode(os.Getenv("GIN_MODE"))

	r := gin.Default()
	r.Use(middlewares.Cors()) // 跨域处理

	// 初始化路由
	initRouters(r)

	return r
}

// 注册http路由
// api/
func initRouters(engine *gin.Engine) {
	engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "health",
		})
	})

	v1 := engine.Group("/api/v1")

	// 验证码 api/v1/captcha
	captchaRouter := v1.Group("/captcha")
	captchaRouter.Use(middlewares.RateLimitMiddleware(2*time.Second, 50))
	{
		captchaRouter.GET("/image", handler.HandleGetImageCode)
		captchaRouter.POST("/sms", handler.HandleGetSmsCode)
	}

	// 用户相关 api/v1/user
	userRouter := v1.Group("/user")
	{
		userRouter.POST("/login", handler.HandleUserLogin)
		userRouter.POST("/login_sms", handler.HandleUserLoginBySms)
		userRouter.POST("/refresh_token", handler.HandleReFreshAccessToken)
		userRouter.POST("/register", handler.HandleUserRegister)
		userRouter.POST("/reset_password", handler.HandleUserResetPassword)
		userRouter.GET("/profile", middlewares.AuthLogin(), handler.HandleGetUserProfile)
	}

	// 题目相关 api/v1/problem
	problemRouter := v1.Group("/problem")
	{
		problemRouter.GET("/tag_list", handler.HandleGetProblemTagList)
		problemRouter.GET("/list", handler.HandleGetProblemList)
		problemRouter.GET("/detail", handler.HandleGetProblemDetail)
		problemRouter.POST("/submit", middlewares.AuthLogin(), handler.HandleSubmitProblem)

		problemRouter.POST("/add", middlewares.AuthLogin(), middlewares.Admin(), handler.HandleCreateProblem)
		problemRouter.POST("/upload_config", middlewares.AuthLogin(), middlewares.Admin(), handler.HandleUploadConfig)
		problemRouter.POST("/publish", middlewares.AuthLogin(), middlewares.Admin(), handler.HandlePublishProblem)
		problemRouter.DELETE("", middlewares.AuthLogin(), middlewares.Admin(), handler.HandleDeleteProblem)
		problemRouter.POST("/update", middlewares.AuthLogin(), middlewares.Admin(), handler.HandleUpdateProblem)
	}

	// record相关 api/v1/record
	recordRouter := v1.Group("/record")
	{
		// 排行榜
		recordRouter.GET("/rank", middlewares.AuthLogin(), handler.HandleGetLeaderboard)
		recordRouter.GET("/result", middlewares.AuthLogin(), handler.HandleGetSubmitResult)        // 本次提交的结果
		recordRouter.GET("/record_list", middlewares.AuthLogin(), handler.HandleGetUserRecordList) // 历史提交记录
		recordRouter.GET("/record", middlewares.AuthLogin(), handler.HandleGetUserRecord)          // 提交记录详情
		recordRouter.GET("/solved_list", middlewares.AuthLogin(), handler.HandleGetUserSolvedList) // 已解决题目
	}

	// 评论 api/v1/comment
	commentRouter := v1.Group("/comment")
	commentRouter.Use(middlewares.AuthLogin())
	{
		commentRouter.POST("/add", handler.HandleCreateComment)
		commentRouter.GET("/root_list", handler.HandleGetRootCommentList)
		commentRouter.GET("/child_list", handler.HandleGetChildCommentList)
		commentRouter.POST("/like", handler.HandleLikeComment)
		commentRouter.DELETE("", handler.HandleDeleteComment)
	}

	// notice api/v1/notice
	noticeRouter := v1.Group("/notice")
	{
		noticeRouter.GET("/list", handler.HandleGetNoticeList)
	}
}
