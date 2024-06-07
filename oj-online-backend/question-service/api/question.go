package api

import (
	"net/http"
	"question-service/api/internal"
	"question-service/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func QuestionSet(ctx *gin.Context) {
	cursor, _ := strconv.Atoi(ctx.DefaultQuery("cursor", "0"))
	if cursor < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		ctx.Abort()
		return
	}
	internal.QuestionSet(ctx, int32(cursor))
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
		internal.QuestionDetail(ctx, idInt)
	}
}

func QuestionQuery(ctx *gin.Context) {
	// 通过题目名字查询相关题目
	if name, ok := ctx.GetQuery("name"); !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
	} else {
		internal.QuestionQuery(ctx, name)
	}
}

func QuestionRun(ctx *gin.Context) {
	// 运行代码
	if form, ok := processOnValidate(ctx, models.QuestionForm{}); ok {
		internal.QuestionRun(ctx, form)
	}
}

func QuestionSubmit(ctx *gin.Context) {
	// 提交代码
	if form, ok := processOnValidate(ctx, models.QuestionForm{}); ok {
		internal.QuestionSubmit(ctx, form)
	}
}
