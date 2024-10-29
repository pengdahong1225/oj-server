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

func (r NoticeHandler) HandleAddNotice(ctx *gin.Context) {
	form, ok := validate(ctx, models.NoticeForm{})
	if !ok {
		return
	}
	uid, _ := ctx.Get("uid")
	res := r.logic.AppendNotice(form, uid.(int64))
	ctx.JSON(http.StatusOK, res)
}

func (r NoticeHandler) HandleDeleteNotice(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "需要告示id",
		})
		ctx.Abort()
		return
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)
	res := r.logic.DeleteNotice(id)
	ctx.JSON(http.StatusOK, res)
}
