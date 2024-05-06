package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"question-service/controller"
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
		engine.GET("/questionSet", controller.QuestionSet)
		engine.POST("/sendSms", middlewares.RateLimitMiddleware(2*time.Second, 50), controller.SendSmsCode)
		// 需要登录
		engine.GET("/rankList", middlewares.AuthLogin(), controller.GetRankList)
	}

	// 题目api
	questionRouter := engine.Group("/questions")
	{
		questionRouter.GET("/detail", controller.QuestionDetail)
		questionRouter.GET("/query", controller.QuestionQuery)
		// 需要登录
		questionRouter.POST("/run", middlewares.AuthLogin(), controller.QuestionRun)
		questionRouter.POST("/submit", middlewares.AuthLogin(), controller.QuestionSubmit)
	}

	// 用户api
	userRouter := engine.Group("/user")
	{
		userRouter.POST("/register", controller.UserRegister)
		userRouter.POST("/login", controller.UserLogin)
		// 需要登录
		userRouter.GET("/detail", middlewares.AuthLogin(), controller.GetUserDetail)
		userRouter.GET("/submitRecord", middlewares.AuthLogin(), controller.GetSubmitRecord)
	}
}

// CmsRouters CMS服务路由
func CmsRouters(engine *gin.Engine) {
	cmsRouter := engine.Group("/cms")
	// 需要管理员权限
	cmsRouter.Use(middlewares.AuthLogin()).Use(middlewares.Admin())
	cmsRouter.GET("/userList", controller.GetUserList)
	cmsRouter.POST("/addQuestion", controller.AddQuestion)
	cmsRouter.POST("/deleteQuestion", controller.DeleteQuestion)
	cmsRouter.POST("/updateQuestion", controller.UpdateQuestion)
}
