package routers

import "github.com/gin-gonic/gin"

// 题目服务路由

func QuestionRouters(engine *gin.Engine) {
	// 题目服务
	questionRouter := engine.Group("/question")
	{
		questionRouter.GET("/getQuestionList", GetQuestionList)
		questionRouter.GET("/getQuestionById", GetQuestionById)
		questionRouter.POST("/addQuestion", AddQuestion)
		questionRouter.POST("/updateQuestion", UpdateQuestion)
		questionRouter.POST("/deleteQuestion", DeleteQuestion)
	}
}
