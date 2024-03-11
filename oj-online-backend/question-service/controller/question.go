package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"question-service/logic"
	"strconv"
)

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

}

func QuestionRun(ctx *gin.Context) {

}

func QuestionSubmit(ctx *gin.Context) {

}
