package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/controller"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/middlewares"
	"net/http"
	"time"
)

// healthCheckRouters 健康检查路由
func healthCheckRouters(engine *gin.Engine) {
	engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    "0",
			"message": "health",
		})
	})
}

// questionRouters
// 题目服务相关路由
// api/
func questionRouters(engine *gin.Engine) {
	// 用户相关
	userRouter := engine.Group("/user")
	{
		userRouter.POST("/login", controller.User{}.HandleLogin)
		userRouter.POST("/register", controller.User{}.HandleRegister)
		userRouter.POST("/reset_password", controller.User{}.HandleResetPassword)
		userRouter.GET("/profile", middlewares.AuthLogin(), controller.User{}.HandleUserProfile)
		userRouter.GET("/record_list", middlewares.AuthLogin(), controller.User{}.HandleRecordList) // 历史提交记录列表
		userRouter.GET("/record", middlewares.AuthLogin(), controller.User{}.HandleRecord)
		userRouter.GET("/solved_list", middlewares.AuthLogin(), controller.User{}.HandleSolvedList)
		userRouter.POST("/refresh_token", controller.User{}.HandleRefreshToken)
	}

	// 题目相关
	problemRouter := engine.Group("/problem")
	{
		problemRouter.GET("/tag_list", controller.ProblemHandler{}.HandleTagList)
		problemRouter.GET("/list", controller.ProblemHandler{}.HandleProblemList)
		problemRouter.GET("/detail", controller.ProblemHandler{}.HandleDetail)
		problemRouter.POST("/submit", middlewares.AuthLogin(), controller.ProblemHandler{}.HandleSubmit)
		problemRouter.GET("/result", middlewares.AuthLogin(), controller.ProblemHandler{}.HandleResult) // 本次提交的结果
		problemRouter.POST("/update", middlewares.AuthLogin(), middlewares.Admin(), controller.ProblemHandler{}.HandleUpdate)
		problemRouter.DELETE("", middlewares.AuthLogin(), middlewares.Admin(), controller.ProblemHandler{}.HandleDelete)
	}

	// 排行榜
	engine.GET("/ranking_list", middlewares.AuthLogin(), controller.User{}.HandleRankList)

	// 评论
	commentRouter := engine.Group("/comment")
	commentRouter.Use(middlewares.AuthLogin())
	{
		commentRouter.POST("/query", controller.CommentHandler{}.HandleGet)
		commentRouter.DELETE("", controller.CommentHandler{}.HandleDelete)
		commentRouter.POST("/add", controller.CommentHandler{}.HandleAdd)
		commentRouter.POST("/like", controller.CommentHandler{}.HandleLike)
	}

	// 验证码
	captchaRouter := engine.Group("/captcha")
	captchaRouter.Use(middlewares.RateLimitMiddleware(2*time.Second, 50))
	{
		captchaRouter.GET("/image", controller.CaptchaHandler{}.GetImageCode)
		captchaRouter.POST("/sms", controller.CaptchaHandler{}.GetSmsCode)
	}

	// 公告
	noticeRouter := engine.Group("/notice")
	{
		noticeRouter.GET("/list", controller.NoticeHandler{}.HandleNoticeList)
		noticeRouter.POST("/add", middlewares.AuthLogin(), middlewares.Admin(), controller.NoticeHandler{}.HandleAddNotice)
		noticeRouter.DELETE("", middlewares.AuthLogin(), middlewares.Admin(), controller.NoticeHandler{}.HandleDeleteNotice)
	}
}
