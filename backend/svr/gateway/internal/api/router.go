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
		user.POST("/logout", uh.HandleUserLogout)
		user.POST("/reset_password", uh.HandleUserResetPassword)
		user.GET("/profile", AuthLogin(), uh.HandleGetUserProfile)
		user.POST("/add", AuthLogin(), Admin(), uh.HandleCreateUser)
		user.GET("/list", AuthLogin(), Admin(), uh.HandleGetUserList)
		user.GET("/captcha/image", RateLimitMiddleware(2*time.Second, 50), uh.HandleGetImageCapatcha)
		user.POST("/captcha/sms", RateLimitMiddleware(2*time.Second, 50), uh.HandleSendSmsCapatcha)
	}

	// 题目
	ph := NewProblemHandler(clients)
	problem := rg.Group("/problem")
	problem.Use(AuthLogin())
	{
		problem.GET("/tag_list", ph.HandleGetProblemTagList)
		problem.GET("/list", ph.HandleGetProblemList)
		problem.GET("/detail", ph.HandleGetProblemDetail)
		problem.POST("/add", Admin(), ph.HandleCreateProblem)
		problem.POST("/publish", Admin(), ph.HandlePublishProblem)
		problem.POST("/hide", Admin(), ph.HandleHideProblem)
		problem.DELETE("", Admin(), ph.HandleDeleteProblem)
		problem.POST("/update", Admin(), ph.HandleUpdateProblem)
		problem.POST("/submit", ph.HandleSubmitProblem)
		problem.GET("/result", ph.HandleGetSubmitResult)
		problem.GET("/solved_list", ph.HandleGetUserSolvedList)
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

	// record（排行榜公开，提交记录需登录）
	rh := NewRecordHandler(clients)
	record := rg.Group("/record")
	{
		record.GET("/rank", rh.HandleGetLeaderboard)                 // 排行榜（公开）
		record.GET("/list", AuthLogin(), rh.HandleGetUserRecordList) // 历史提交记录列表
		record.GET("/detail", AuthLogin(), rh.HandleGetUserRecord)   // 提交记录详情
	}

	// notice（公告 list/detail 公开，增删改/发布/隐藏 需管理员）
	nh := NewNoticeHandler(clients)
	notice := rg.Group("/notice")
	{
		notice.GET("/list", nh.HandleGetNoticeList)     // 公告列表（公开）
		notice.GET("/detail", nh.HandleGetNoticeDetail) // 公告详情（公开）
		notice.POST("/add", AuthLogin(), Admin(), nh.HandleAddNotice)
		notice.DELETE("/delete", AuthLogin(), Admin(), nh.HandleDeleteNotice)
		notice.POST("/update", AuthLogin(), Admin(), nh.HandleUpdateNotice)
		notice.POST("/publish", AuthLogin(), Admin(), nh.HandlePublishNotice)
		notice.POST("/hide", AuthLogin(), Admin(), nh.HandleHideNotice)
	}
}
