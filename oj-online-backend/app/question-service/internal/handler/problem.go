package handler

import (
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/logic"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/middlewares"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProblemHandler struct{}

// ProblemSet 题目列表
// 查询参数 cursor：游标，代表下一个元素的id
func (receiver ProblemHandler) ProblemSet(ctx *gin.Context) {
	cursor, _ := strconv.Atoi(ctx.DefaultQuery("cursor", "0"))
	if cursor < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
		})
		ctx.Abort()
		return
	}
	res := logic.ProblemLogic{}.HandleProblemSet(cursor)
	ctx.JSON(res.Code, res)
}

func (receiver ProblemHandler) ProblemDetail(ctx *gin.Context) {
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
	res := logic.ProblemLogic{}.HandleProblemDetail(ID)
	ctx.JSON(res.Code, res)
}

func (receiver ProblemHandler) ProblemSubmit(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middlewares.UserClaims)
	// 表单验证
	form, ret := validate(ctx, models.SubmitForm{})
	if !ret {
		return
	}

	res := logic.ProblemLogic{}.HandleProblemSubmit(claims.Uid, form)
	ctx.JSON(res.Code, res)
}

// 查询提交结果
func (receiver ProblemHandler) QueryResult(ctx *gin.Context) {
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
	res := logic.ProblemLogic{}.HandleQueryResult(claims.Uid, problemID)
	ctx.JSON(res.Code, res)
}

func (receiver ProblemHandler) ProblemSearch(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
		})
		return
	}
	res := logic.ProblemLogic{}.HandleProblemSearch(name)
	ctx.JSON(res.Code, res)
}
