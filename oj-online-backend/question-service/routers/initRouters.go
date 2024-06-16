package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"question-service/api"
	"question-service/middlewares"
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

// QuestionRouters 服务路由
func QuestionRouters(engine *gin.Engine) {
	// api/
	{
		engine.GET("/problemSet", api.ProblemSet)
		engine.GET("/rankList", middlewares.AuthLogin(), api.GetRankList)
		engine.POST("/login", api.Login)
	}

	// api/user
	userRouter := engine.Group("/user")
	{
		userRouter.Use(middlewares.AuthLogin())
		userRouter.GET("/profile", api.GetUserProfile)
		userRouter.GET("/submitRecord", api.GetSubmitRecord)
		userRouter.GET("/solvedList", api.GetUserSolvedList)
	}

	// api/captcha
	captchaRouter := engine.Group("/captcha")
	{
		captchaRouter.Use(middlewares.RateLimitMiddleware(2*time.Second, 50))
		captchaRouter.GET("/image", api.GetImageCode)
		captchaRouter.POST("/sms", api.GetSmsCode)
	}

	// api/problem
	problemRouter := engine.Group("/problem")
	{
		problemRouter.GET("/detail", api.GetProblemDetail)
		problemRouter.GET("/search", api.QueryProblem)
		// 需要登录
		problemRouter.POST("/submit", middlewares.AuthLogin(), api.ProblemSubmit)
	}
}

// CmsRouters cms路由
func CmsRouters(engine *gin.Engine) {
	cmsRouter := engine.Group("/cms")
	// 需要管理员权限
	cmsRouter.Use(middlewares.AuthLogin()).Use(middlewares.Admin())
	cmsRouter.GET("/userList", api.GetUserList)
	cmsRouter.POST("/addQuestion", api.AddQuestion)
	cmsRouter.POST("/deleteQuestion", api.DeleteQuestion)
	cmsRouter.POST("/updateQuestion", api.UpdateQuestion)
}
