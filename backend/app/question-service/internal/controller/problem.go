package controller

import (
	"encoding/json"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/logic"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/middlewares"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProblemHandler struct {
	logic logic.ProblemLogic
}

// HandleProblemSet 题目列表
func (r ProblemHandler) HandleProblemSet(ctx *gin.Context) {
	p := ctx.Query("params")
	if p == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "参数错误",
		})
		ctx.Abort()
		return
	}
	params := &models.QueryProblemListParams{}
	err := json.Unmarshal([]byte(p), params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": err.Error(),
		})
		ctx.Abort()
		return
	}
	if params.Page <= 0 || params.PageSize <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "参数错误",
		})
		ctx.Abort()
		return
	}

	res := r.logic.GetProblemList(params)
	ctx.JSON(http.StatusOK, res)
}

func (r ProblemHandler) HandleDetail(ctx *gin.Context) {
	// 查询参数
	idStr := ctx.Query("problemID")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "参数错误",
		})
		return
	}

	id, _ := strconv.ParseInt(idStr, 10, 64)
	res := r.logic.GetProblemDetail(id)
	ctx.JSON(http.StatusOK, res)
}

func (r ProblemHandler) HandleSubmit(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middlewares.UserClaims)
	// 表单验证
	form, ret := validate(ctx, models.SubmitForm{})
	if !ret {
		return
	}

	res := r.logic.OnProblemSubmit(claims.Uid, form)
	ctx.JSON(http.StatusOK, res)
}

// 查询提交结果
func (r ProblemHandler) HandleResult(ctx *gin.Context) {
	// 查询参数
	id := ctx.Query("problemID")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "参数错误",
		})
		return
	}
	problemID, _ := strconv.ParseInt(id, 10, 64)
	claims := ctx.MustGet("claims").(*middlewares.UserClaims)
	res := r.logic.QueryResult(claims.Uid, problemID)
	ctx.JSON(http.StatusOK, res)
}
