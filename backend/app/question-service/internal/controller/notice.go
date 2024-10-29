package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/logic"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/models"
	"net/http"
)

type NoticeHandler struct {
	logic logic.NoticeLogic
}

func (r NoticeHandler) HandleNoticeList(ctx *gin.Context) {
	p := ctx.Query("params")
	if p == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "参数错误",
		})
		ctx.Abort()
		return
	}
	params := &models.QueryNoticeListParams{}
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

	res := r.logic.GetNoticeList(params)
	ctx.JSON(http.StatusOK, res)
}

func (r NoticeHandler) HandleAddNotice(ctx *gin.Context) {}

func (r NoticeHandler) HandleDeleteNotice(ctx *gin.Context) {}
