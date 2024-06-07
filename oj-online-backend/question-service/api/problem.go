package api

import (
	"net/http"
	"question-service/api/internal"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProblemSet 题目列表
// 查询参数 cursor：游标，代表下一个元素的id
func ProblemSet(ctx *gin.Context) {
	cursor, _ := strconv.Atoi(ctx.DefaultQuery("cursor", "0"))
	if cursor < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
		})
		ctx.Abort()
		return
	}
	res := internal.ProblemSet(cursor)
	ctx.JSON(res.Code, res)
}

func QuestionDetail(ctx *gin.Context) {

}

func QuestionQuery(ctx *gin.Context) {

}

func QuestionRun(ctx *gin.Context) {

}

func QuestionSubmit(ctx *gin.Context) {

}
