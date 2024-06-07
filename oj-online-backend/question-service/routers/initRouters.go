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
			"msg": "health",
		})
	})
}

// QuestionRouters 题目服务路由
func QuestionRouters(engine *gin.Engine) {
	// 公共api
	{
		engine.GET("/questionSet", api.QuestionSet)
		engine.POST("/sendSms", middlewares.RateLimitMiddleware(2*time.Second, 50), api.SendSmsCode)
		// 需要登录
		engine.GET("/rankList", middlewares.AuthLogin(), api.GetRankList)
	}

	// 题目api
	questionRouter := engine.Group("/questions")
	{
		questionRouter.GET("/detail", api.QuestionDetail)
		questionRouter.GET("/query", api.QuestionQuery)
		// 需要登录
		questionRouter.POST("/run", middlewares.AuthLogin(), api.QuestionRun)
		questionRouter.POST("/submit", middlewares.AuthLogin(), api.QuestionSubmit)
	}

	// 用户api
	userRouter := engine.Group("/user")
	{
		userRouter.POST("/login", api.UserLogin)
		// 需要登录
		userRouter.GET("/detail", middlewares.AuthLogin(), api.GetUserDetail)
		userRouter.GET("/submitRecord", middlewares.AuthLogin(), api.GetSubmitRecord)
	}
}

// CmsRouters CMS服务路由
func CmsRouters(engine *gin.Engine) {
	cmsRouter := engine.Group("/cms")
	// 需要管理员权限
	cmsRouter.Use(middlewares.AuthLogin()).Use(middlewares.Admin())
	cmsRouter.GET("/userList", api.GetUserList)
	cmsRouter.POST("/addQuestion", api.AddQuestion)
	cmsRouter.POST("/deleteQuestion", api.DeleteQuestion)
	cmsRouter.POST("/updateQuestion", api.UpdateQuestion)
}
