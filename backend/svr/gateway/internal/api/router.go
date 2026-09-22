package api

import (
	"net/http"
	"oj-server/svr/gateway/internal/client"
	"time"

	"github.com/gin-gonic/gin"
)

func Router(clients *client.Clients) *gin.Engine {
	gin.SetMode(gin.DebugMode)

	engine := gin.New()
	engine.Use(Logger(), Recovery())
	engine.Use(Cors())

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "health",
		})
	})

	v1 := engine.Group("/api/v1")

	initRouters(v1, clients)

	return engine
}
func initRouters(rg *gin.RouterGroup, clients *client.Clients) {
	// 用户
	uh := NewUserHandler(clients)
	user := rg.Group("/user")
	{
		user.POST("/login", uh.HandleUserLogin)
		user.POST("/login_sms", uh.HandleUserLoginBySms)
		user.POST("/refresh_token", uh.HandleReFreshAccessToken)
		user.POST("/register", uh.HandleUserRegister)
		user.POST("/reset_password", uh.HandleUserResetPassword)
		user.GET("/profile", AuthLogin(), uh.HandleGetUserProfile)
		user.GET("/list", AuthLogin(), Admin(), uh.HandleGetUserList)
	}

	// 题目
	ph := NewProblemHandler(clients)
	problem := rg.Group("/problem")
	problem.Use(AuthLogin())
	{
		problem.POST("/add", Admin(), ph.HandleCreateProblem)
		problem.POST("/upload_config", Admin(), ph.HandleUploadConfig)
		problem.GET("/download_config", Admin(), ph.HandleDownloadConfig)
		problem.POST("/publish", Admin(), ph.HandlePublishProblem)
		problem.POST("/hide", Admin(), ph.HandleHideProblem)
		problem.DELETE("", Admin(), ph.HandleDeleteProblem)
		problem.POST("/update", Admin(), ph.HandleUpdateProblem)
		problem.GET("/tag_list", ph.HandleGetProblemTagList)
		problem.GET("/list", ph.HandleGetProblemList)
		problem.GET("/detail", ph.HandleGetProblemDetail)
		problem.POST("/submit", ph.HandleSubmitProblem)
	}

	// 评论
	ch := NewCommentHandler(clients)
	comment := rg.Group("/comment")
	comment.Use(AuthLogin())
	{
		comment.POST("/add", ch.HandleCreateComment)
		comment.GET("/root_list", ch.HandleGetRootCommentList)
		comment.GET("/child_list", ch.HandleGetChildCommentList)
		comment.POST("/like", ch.HandleLikeComment)
		comment.DELETE("", ch.HandleDeleteComment)
	}

	// record
	rh := NewRecordHandler(clients)
	record := rg.Group("/record")
	record.Use(AuthLogin())
	{
		record.GET("/rank", rh.HandleGetLeaderboard)           // 排行榜
		record.GET("/result", rh.HandleGetSubmitResult)        // 本次提交的结果
		record.GET("/record_list", rh.HandleGetUserRecordList) // 历史提交记录
		record.GET("/record", rh.HandleGetUserRecord)          // 提交记录详情
		record.GET("/solved_list", rh.HandleGetUserSolvedList) // 已解决题目
	}

	// 验证码
	captcha := rg.Group("/captcha")
	captcha.Use(RateLimitMiddleware(2*time.Second, 50))
	{
		captcha.GET("/image", HandleGetImageCode)
		captcha.POST("/sms", HandleGetSmsCode)
	}
}
