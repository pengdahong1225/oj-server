package router

import (
	"github.com/gin-gonic/gin"
	"oj-server/svr/gateway/internal/handler"
	"oj-server/svr/gateway/internal/middlewares"
)

func initProblemRouter(rg *gin.RouterGroup) {
	// 题目相关
	problem := rg.Group("/problem")
	{
		problem.GET("/tag_list", handler.HandleGetProblemTagList)
		problem.GET("/list", handler.HandleGetProblemList)
		problem.GET("/detail", handler.HandleGetProblemDetail)
		problem.POST("/submit", middlewares.AuthLogin(), handler.HandleSubmitProblem)

		problem.POST("/add", middlewares.AuthLogin(), middlewares.Admin(), handler.HandleCreateProblem)
		problem.POST("/upload_config", middlewares.AuthLogin(), middlewares.Admin(), handler.HandleUploadConfig)
		problem.POST("/publish", middlewares.AuthLogin(), middlewares.Admin(), handler.HandlePublishProblem)
		problem.DELETE("", middlewares.AuthLogin(), middlewares.Admin(), handler.HandleDeleteProblem)
		problem.POST("/update", middlewares.AuthLogin(), middlewares.Admin(), handler.HandleUpdateProblem)
	}

	// 评论相关
	comment := rg.Group("/comment")
	comment.Use(middlewares.AuthLogin())
	{
		comment.POST("/add", handler.HandleCreateComment)
		comment.GET("/root_list", handler.HandleGetRootCommentList)
		comment.GET("/child_list", handler.HandleGetChildCommentList)
		comment.POST("/like", handler.HandleLikeComment)
		comment.DELETE("", handler.HandleDeleteComment)
	}

	// record 相关
	record := rg.Group("/record")
	record.Use(middlewares.AuthLogin())
	{
		// 排行榜
		record.GET("/rank", handler.HandleGetLeaderboard)
		record.GET("/result", handler.HandleGetSubmitResult)        // 本次提交的结果
		record.GET("/record_list", handler.HandleGetUserRecordList) // 历史提交记录
		record.GET("/record", handler.HandleGetUserRecord)          // 提交记录详情
		record.GET("/solved_list", handler.HandleGetUserSolvedList) // 已解决题目
	}
}
