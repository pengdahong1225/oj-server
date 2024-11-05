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
func questionRouters(engine *gin.Engine) {
	engine.Use(middlewares.CheckIsExpiration())
	// api/
	{
		engine.GET("/problemSet", controller.ProblemHandler{}.HandleProblemSet)
		engine.GET("/rankList", middlewares.AuthLogin(), controller.User{}.HandleRankList)
		engine.POST("/login", controller.User{}.HandleLogin)
		engine.POST("/register", controller.User{}.HandleRegister)
		engine.POST("/resetPassword", controller.User{}.HandleResetPassword)
	}

	// api/user
	userRouter := engine.Group("/user")
	{
		userRouter.Use(middlewares.AuthLogin())
		userRouter.GET("/profile", controller.User{}.HandleUserProfile)
		userRouter.GET("/submitRecord", controller.User{}.HandleSubmitRecord) // 历史提交记录
		userRouter.GET("/solvedList", controller.User{}.HandleSolvedList)
	}

	// api/problem
	problemRouter := engine.Group("/problem")
	{
		problemRouter.GET("/detail", controller.ProblemHandler{}.HandleDetail)
		problemRouter.POST("/submit", middlewares.AuthLogin(), controller.ProblemHandler{}.HandleSubmit)
		problemRouter.GET("/result", middlewares.AuthLogin(), controller.ProblemHandler{}.HandleResult) // 本次提交的结果
		problemRouter.POST("/update", middlewares.AuthLogin(), middlewares.Admin(), controller.ProblemHandler{}.HandleUpdate)
		problemRouter.POST("/delete", middlewares.AuthLogin(), middlewares.Admin(), controller.ProblemHandler{}.HandleDelete)
		problemRouter.GET("/tagList", controller.ProblemHandler{}.HandleTagList)
	}

	// api/comment
	commentRouter := engine.Group("/comment")
	commentRouter.Use(middlewares.AuthLogin())
	{
		commentRouter.POST("/add", controller.CommentHandler{}.HandleAdd)
		commentRouter.POST("/get", controller.CommentHandler{}.HandleGet)
		commentRouter.POST("/delete", controller.CommentHandler{}.HandleDelete)
		commentRouter.POST("/like", controller.CommentHandler{}.HandleLike)
	}

	// api/captcha
	captchaRouter := engine.Group("/captcha")
	captchaRouter.Use(middlewares.RateLimitMiddleware(2*time.Second, 50))
	{
		captchaRouter.GET("/image", controller.CaptchaHandler{}.GetImageCode)
		captchaRouter.POST("/sms", controller.CaptchaHandler{}.GetSmsCode)
	}

	// api/notice
	noticeRouter := engine.Group("/notice")
	{
		noticeRouter.GET("/noticeList", controller.NoticeHandler{}.HandleNoticeList)
		noticeRouter.POST("/addNotice", middlewares.AuthLogin(), middlewares.Admin(), controller.NoticeHandler{}.HandleAddNotice)
		noticeRouter.DELETE("/delNotice", middlewares.AuthLogin(), middlewares.Admin(), controller.NoticeHandler{}.HandleDeleteNotice)
	}
}
