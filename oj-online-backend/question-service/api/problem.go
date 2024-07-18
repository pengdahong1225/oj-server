package api

import (
	"net/http"
	"question-service/api/internal"
	"question-service/middlewares"
	"question-service/models"
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
	res := internal.ProblemHandler{}.GetProblemSet(cursor)
	ctx.JSON(res.Code, res)
}

func ProblemDetail(ctx *gin.Context) {
	// 查询参数
	problemID := ctx.Query("problemID")
	if problemID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
		})
		return
	}

	ID, _ := strconv.ParseInt(problemID, 10, 64)
	res := internal.ProblemHandler{}.GetProblemDetail(ID)
	ctx.JSON(res.Code, res)
}

func ProblemSearch(ctx *gin.Context) {}

func ProblemSubmit(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middlewares.UserClaims)
	// 表单验证
	form, ret := validate(ctx, models.SubmitForm{})
	if !ret {
		return
	}

	res := internal.ProblemHandler{}.ProblemSubmit(claims.Uid, form)
	ctx.JSON(res.Code, res)
}

func QueryResult(ctx *gin.Context) {

}
