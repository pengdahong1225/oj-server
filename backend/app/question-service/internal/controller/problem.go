package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/logic"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/models"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/svc/cache"
	"net/http"
	"strconv"
)

type ProblemHandler struct {
	logic logic.ProblemLogic
}

func (r ProblemHandler) HandleUpdate(ctx *gin.Context) {
	uid := ctx.GetInt64("uid")
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

	res := r.logic.UpdateQuestion(uid, form, &config)
	ctx.JSON(http.StatusOK, res)
}

func (r ProblemHandler) HandleDelete(ctx *gin.Context) {
	p := ctx.Query("problem_id")
	if p == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "参数错误",
		})
		return
	}
	problemID, _ := strconv.ParseInt(p, 10, 64)
	res := r.logic.DeleteProblem(problemID)
	ctx.JSON(http.StatusOK, res)
}

// HandleProblemList
// 题目列表
func (r ProblemHandler) HandleProblemList(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("page_size")
	keyWord := ctx.Query("keyword")
	tag := ctx.Query("tag")
	page, err := strconv.Atoi(pageStr)
	uidStr := ctx.DefaultQuery("uid", "")
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

	uid, _ := strconv.ParseInt(uidStr, 10, 64)
	res := r.logic.GetProblemList(params, uid)
	ctx.JSON(http.StatusOK, res)
}

func (r ProblemHandler) HandleDetail(ctx *gin.Context) {
	// 查询参数
	idStr := ctx.Query("problem_id")
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
	uid := ctx.GetInt64("uid")
	// 表单验证
	form, ret := validate(ctx, models.SubmitForm{})
	if !ret {
		return
	}

	res := r.logic.OnProblemSubmit(uid, form)
	ctx.JSON(http.StatusOK, res)
}

// HandleResult 查询提交结果
// @problemID 题目ID
func (r ProblemHandler) HandleResult(ctx *gin.Context) {
	// 查询参数
	idStr := ctx.Query("problem_id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "参数错误",
		})
		return
	}
	problemID, _ := strconv.ParseInt(idStr, 10, 64)
	uid := ctx.GetInt64("uid")

	res := r.logic.QueryResult(uid, problemID)
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
