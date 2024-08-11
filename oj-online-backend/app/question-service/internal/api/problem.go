package api

import (
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/api/handler"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/middlewares"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/models"
	"net/http"
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
	res := handler.ProblemHandler{}.HandleProblemSet(cursor)
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
	res := handler.ProblemHandler{}.HandleProblemDetail(ID)
	ctx.JSON(res.Code, res)
}

func ProblemSubmit(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middlewares.UserClaims)
	// 表单验证
	form, ret := validate(ctx, models.SubmitForm{})
	if !ret {
		return
	}

	res := handler.ProblemHandler{}.HandleProblemSubmit(claims.Uid, form)
	ctx.JSON(res.Code, res)
}

// 查询提交结果
func QueryResult(ctx *gin.Context) {
	// 查询参数
	id := ctx.Query("problemID")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
		})
		return
	}
	problemID, _ := strconv.ParseInt(id, 10, 64)
	claims := ctx.MustGet("claims").(*middlewares.UserClaims)
	res := handler.ProblemHandler{}.HandleQueryResult(claims.Uid, problemID)
	ctx.JSON(res.Code, res)
}

func ProblemSearch(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
		})
		return
	}
	res := handler.ProblemHandler{}.HandleProblemSearch(name)
	ctx.JSON(res.Code, res)
}
