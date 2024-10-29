package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/logic"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/models"
	"net/http"
	"strconv"
)

type NoticeHandler struct {
	logic logic.NoticeLogic
}

func (r NoticeHandler) HandleNoticeList(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("page_size")
	keyWord := ctx.Query("keyword")
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

	params := &models.QueryNoticeListParams{
		Page:     int32(page),
		PageSize: int32(pageSize),
		KeyWord:  keyWord,
	}
	res := r.logic.GetNoticeList(params)
	ctx.JSON(http.StatusOK, res)
}

func (r NoticeHandler) HandleAddNotice(ctx *gin.Context) {}

func (r NoticeHandler) HandleDeleteNotice(ctx *gin.Context) {}
