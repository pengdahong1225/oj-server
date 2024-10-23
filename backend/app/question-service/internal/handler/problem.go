package handler

import (
	"encoding/json"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/logic"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/middlewares"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProblemHandler struct{}

// HandleProblemSet 题目列表
func (receiver ProblemHandler) HandleProblemSet(ctx *gin.Context) {
	p := ctx.Query("params")
	if p == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
		})
		ctx.Abort()
		return
	}
	params := &models.QueryProblemListParams{}
	err := json.Unmarshal([]byte(p), params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		ctx.Abort()
		return
	}
	if params.Page <= 0 || params.PageSize <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
		})
		ctx.Abort()
		return
	}

	res := logic.ProblemLogic{}.GetProblemList(params)
	ctx.JSON(res.Code, res)
}

func (receiver ProblemHandler) HandleDetail(ctx *gin.Context) {
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

func (receiver ProblemHandler) HandleSubmit(ctx *gin.Context) {
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
func (receiver ProblemHandler) HandleResult(ctx *gin.Context) {
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

func (receiver ProblemHandler) HandleSearch(ctx *gin.Context) {
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
