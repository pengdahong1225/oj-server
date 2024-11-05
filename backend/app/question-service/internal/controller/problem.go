package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/logic"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/middlewares"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/models"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/svc/cache"
	"net/http"
	"strconv"
)

type ProblemHandler struct {
	logic logic.ProblemLogic
}

func (r ProblemHandler) HandleUpdate(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middlewares.UserClaims)
	form, ret := validate(ctx, models.UpdateProblemForm{})
	if !ret {
		return
	}

	config := models.ProblemConfig{}
	err := json.Unmarshal([]byte(form.Config), &config)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "config json 解析失败，格式错误",
		})
		return
	}

	res := r.logic.UpdateQuestion(claims.Uid, form, &config)
	ctx.JSON(http.StatusOK, res)
}

func (r ProblemHandler) HandleDelete(ctx *gin.Context) {
	p := ctx.GetString("problem_id")
	if p == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "参数错误",
		})
		return
	}
	problemID, _ := strconv.ParseInt(p, 10, 64)
	res := r.logic.DeleteQuestion(problemID)
	ctx.JSON(http.StatusOK, res)
}

// HandleProblemSet 题目列表
func (r ProblemHandler) HandleProblemSet(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("page_size")
	keyWord := ctx.Query("keyword")
	tag := ctx.Query("tag")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "页码参数错误",
		})
		ctx.Abort()
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "页大小参数错误",
		})
		ctx.Abort()
		return
	}
	params := &models.QueryProblemListParams{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Keyword:  keyWord,
		Tag:      tag,
	}

	c, _ := ctx.Get("claims")
	claim := c.(*middlewares.UserClaims)

	res := r.logic.GetProblemList(params, claim.Uid)
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

// HandleResult 查询提交结果
// @problemID 题目ID
func (r ProblemHandler) HandleResult(ctx *gin.Context) {
	// 查询参数
	idStr := ctx.Query("problemID")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "参数错误",
		})
		return
	}
	problemID, _ := strconv.ParseInt(idStr, 10, 64)
	claims := ctx.MustGet("claims").(*middlewares.UserClaims)

	res := r.logic.QueryResult(claims.Uid, problemID)
	ctx.JSON(http.StatusOK, res)
}

func (r ProblemHandler) HandleTagList(ctx *gin.Context) {
	tags, err := cache.GetTagList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    models.Failed,
			"message": fmt.Sprintf("获取标签列表失败:%s", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    models.Success,
		"message": "OK",
		"data":    tags,
	})
}
