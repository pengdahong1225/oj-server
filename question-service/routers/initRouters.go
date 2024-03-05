package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"question-service/controller"
	"question-service/middlewares"
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
	// 公共方法
	engine.GET("/questionSet", controller.QuestionSet)
	engine.POST("/sendCms", controller.SendCmsCode)

	// 题目
	questionRouter := engine.Group("/questions")
	{
		questionRouter.GET("/detail", controller.QuestionDetail)
		questionRouter.GET("/query", controller.QuestionQuery)
		// 需要登录
		questionRouter.POST("/run", middlewares.AuthLogin(), controller.QuestionRun)
		questionRouter.POST("/submit", middlewares.AuthLogin(), controller.QuestionSubmit)
	}

	// 用户
	userRouter := engine.Group("/user")
	{
		userRouter.POST("/register", controller.UserRegister)
		userRouter.POST("/login", controller.UserLogin)
		// 需要登录
		userRouter.GET("/detail", middlewares.AuthLogin(), controller.GetUserDetail)
		userRouter.GET("/rankList", middlewares.AuthLogin(), controller.GetRankList)
		userRouter.GET("/submitRecord", middlewares.AuthLogin(), controller.GetSubmitRecord)
	}
}
