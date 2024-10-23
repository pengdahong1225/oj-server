package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/handler"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/middlewares"
	"net/http"
	"time"
)

// HealthCheckRouters 健康检查路由
func HealthCheckRouters(engine *gin.Engine) {
	engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    "ok",
			"message": "health",
		})
	})
}

func QuestionRouters(engine *gin.Engine) {
	// api/
	{
		engine.GET("/problemSet", handler.ProblemHandler{}.HandleProblemSet)
		engine.GET("/rankList", middlewares.AuthLogin(), handler.User{}.HandleRankList)
		engine.POST("/login", handler.User{}.HandleLogin)
		engine.POST("/register", handler.User{}.HandleRegister)
		engine.POST("/resetPassword", handler.User{}.HandleResetPassword)
	}

	// api/user
	userRouter := engine.Group("/user")
	{
		userRouter.Use(middlewares.AuthLogin())
		userRouter.GET("/profile", handler.User{}.HandleUserProfile)
		userRouter.GET("/submitRecord", handler.User{}.HandleSubmitRecord)
		userRouter.GET("/solvedList", handler.User{}.HandleSolvedList)
	}

	// api/problem
	problemRouter := engine.Group("/problem")
	{
		problemRouter.GET("/detail", handler.ProblemHandler{}.HandleDetail)
		problemRouter.GET("/search", handler.ProblemHandler{}.HandleSearch)
		// 需要登录
		problemRouter.POST("/submit", middlewares.AuthLogin(), handler.ProblemHandler{}.HandleSubmit)
		problemRouter.GET("/result", middlewares.AuthLogin(), handler.ProblemHandler{}.HandleResult)
	}

	// api/comment
	commentRouter := engine.Group("/comment")
	//commentRouter.Use(middlewares.AuthLogin())
	{
		commentRouter.POST("/add", handler.CommentHandler{}.HandleAdd)
		commentRouter.POST("/get", handler.CommentHandler{}.HandleGet)
		commentRouter.POST("/delete", handler.CommentHandler{}.HandleDelete)
		commentRouter.POST("/like", handler.CommentHandler{}.HandleLike)
	}

	// api/captcha
	captchaRouter := engine.Group("/captcha")
	captchaRouter.Use(middlewares.RateLimitMiddleware(2*time.Second, 50))
	{
		captchaRouter.GET("/image", handler.CaptchaHandler{}.GetImageCode)
		captchaRouter.POST("/sms", handler.CaptchaHandler{}.GetSmsCode)
	}
}

// CmsRouters cms路由
// 题目的增删改操作
func CmsRouters(engine *gin.Engine) {
	cmsRouter := engine.Group("/cms")
	// 需要管理员权限
	cmsRouter.Use(middlewares.AuthLogin()).Use(middlewares.Admin())
	cmsRouter.POST("/updateQuestion", handler.AdminHandler{}.UpdateQuestion)
	cmsRouter.POST("/deleteQuestion", handler.AdminHandler{}.DeleteQuestion)
}
