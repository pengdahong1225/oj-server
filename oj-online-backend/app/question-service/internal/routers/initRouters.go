package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/handler"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/middlewares"
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
		engine.GET("/problemSet", handler.ProblemHandler{}.ProblemSet)
		engine.GET("/rankList", middlewares.AuthLogin(), handler.UserHandler{}.GetRankList)
		engine.POST("/login", handler.UserHandler{}.Login)
	}

	// api/user
	userRouter := engine.Group("/user")
	{
		userRouter.Use(middlewares.AuthLogin())
		userRouter.GET("/profile", handler.UserHandler{}.GetUserProfile)
		userRouter.GET("/submitRecord", handler.UserHandler{}.GetSubmitRecord)
		userRouter.GET("/solvedList", handler.UserHandler{}.GetUserSolvedList)
	}

	// api/captcha
	captchaRouter := engine.Group("/captcha")
	captchaRouter.Use(middlewares.RateLimitMiddleware(2*time.Second, 50))
	{
		captchaRouter.GET("/image", handler.CaptchaHandler{}.GetImageCode)
		captchaRouter.POST("/sms", handler.CaptchaHandler{}.GetSmsCode)
	}

	// api/problem
	problemRouter := engine.Group("/problem")
	{
		problemRouter.GET("/detail", handler.ProblemHandler{}.ProblemDetail)
		problemRouter.GET("/search", handler.ProblemHandler{}.ProblemSearch)
		// 需要登录
		problemRouter.POST("/submit", middlewares.AuthLogin(), handler.ProblemHandler{}.ProblemSubmit)
		problemRouter.GET("/queryResult", middlewares.AuthLogin(), handler.ProblemHandler{}.QueryResult)
	}

	// api/comment
	commentRouter := engine.Group("/comment")
	commentRouter.Use(middlewares.AuthLogin())
	{
		commentRouter.POST("/add", handler.CommentHandler{}.Add)
		commentRouter.POST("/delete", handler.CommentHandler{}.Delete)
		commentRouter.POST("/like", handler.CommentHandler{}.Like)
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
