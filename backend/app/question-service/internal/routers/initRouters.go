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
			"code":    "ok",
			"message": "health",
		})
	})
}

// questionRouters
// 题目服务相关路由
func questionRouters(engine *gin.Engine) {
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
		userRouter.GET("/submitRecord", controller.User{}.HandleSubmitRecord)
		userRouter.GET("/upss", controller.User{}.HandleUPSS)
		userRouter.GET("/solvedList", controller.User{}.HandleSolvedList)
	}

	// api/problem
	problemRouter := engine.Group("/problem")
	{
		problemRouter.GET("/detail", controller.ProblemHandler{}.HandleDetail)
		// 需要登录
		problemRouter.POST("/submit", middlewares.AuthLogin(), controller.ProblemHandler{}.HandleSubmit)
		problemRouter.GET("/result", middlewares.AuthLogin(), controller.ProblemHandler{}.HandleResult)
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
}

// cmsRouters
// cms路由
func cmsRouters(engine *gin.Engine) {
	cmsRouter := engine.Group("/cms")
	// 需要管理员权限
	cmsRouter.Use(middlewares.AuthLogin()).Use(middlewares.Admin())
	cmsRouter.POST("/updateQuestion", controller.AdminHandler{}.UpdateQuestion)
	cmsRouter.POST("/deleteQuestion", controller.AdminHandler{}.DeleteQuestion)
}

// noticeRouters
// notice路由
func noticeRouters(engine *gin.Engine) {
	noticeRouter := engine.Group("/notice")
	noticeRouter.GET("/noticeList", controller.NoticeHandler{}.HandleNoticeList)
	noticeRouter.POST("/addNotice", middlewares.Admin(), controller.NoticeHandler{}.HandleAddNotice)
	noticeRouter.POST("/deleteNotice", middlewares.Admin(), controller.NoticeHandler{}.HandleDeleteNotice)
}
