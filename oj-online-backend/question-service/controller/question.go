package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"question-service/logic"
	"question-service/models"
	"strconv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func QuestionSet(ctx *gin.Context) {
	// 题库分页
	cursor, _ := strconv.Atoi(ctx.DefaultQuery("cursor", "0"))
	if cursor < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		ctx.Abort()
		return
	}
	logic.QuestionSet(ctx, int32(cursor))
}

func QuestionDetail(ctx *gin.Context) {
	// 题目详情
	if id, ok := ctx.GetQuery("id"); !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		ctx.Abort()
		return
	} else {
		idInt, _ := strconv.ParseInt(id, 10, 64)
		logic.QuestionDetail(ctx, idInt)
	}
}

func QuestionQuery(ctx *gin.Context) {
	// 通过题目名字查询相关题目
	if name, ok := ctx.GetQuery("name"); !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
	} else {
		logic.QuestionQuery(ctx, name)
	}
}

func QuestionRun(ctx *gin.Context) {
	// 运行代码
	if form, ok := processOnValidate(ctx, models.QuestionForm{}); ok {
		// 升级连接
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "连接失败",
			})
			ctx.Abort()
			return
		}
		logic.QuestionRun(ctx, form, conn)
	}
}

func QuestionSubmit(ctx *gin.Context) {
	// 提交代码
	if form, ok := processOnValidate(ctx, models.QuestionForm{}); ok {
		// 升级连接
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "连接失败",
			})
			ctx.Abort()
			return
		}
		logic.QuestionSubmit(ctx, form, conn)
	}
}

func JudgeCallback(ctx *gin.Context) {
	if form, ok := processOnValidate(ctx, models.JudgeBackForm{}); ok {
		logic.JudgeCallback(ctx, form)
	}
}
