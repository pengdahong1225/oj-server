package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/api/IpQuery"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/logic"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

type CommentHandler struct {
	logic logic.CommentLogic
}

func (r CommentHandler) HandleAdd(ctx *gin.Context) {
	form, ok := validate(ctx, models.AddCommentForm{})
	if !ok {
		return
	}

	// 查询ip归属地
	ip := ctx.ClientIP()
	info, err := IpQuery.QueryIpGeolocation(ip)
	if err != nil {
		logrus.Errorf("查询ip归属地失败,ip:%s,err:%s", ip, err.Error())
		info.RegionName = "未知地区"
	}

	res := r.logic.OnAddComment(form, info.RegionName)
	ctx.JSON(http.StatusOK, res)
}

// HandleGetRootList 顶层评论查询
func (r CommentHandler) HandleGetRootList(ctx *gin.Context) {
	// 查询参数校验
	params := &models.RootCommentListQueryParams{}
	err := ctx.ShouldBindQuery(params)
	if err != nil {
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    models.Failed,
				"message": "表单验证错误",
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": errs.Error(),
		})
		return
	}

	res := r.logic.OnQueryRootCommentList(params)
	ctx.JSON(http.StatusOK, res)
}

func (r CommentHandler) HandleGetChildList(ctx *gin.Context) {
	// 查询参数校验
	params := &models.ChildCommentListQueryParams{}
	err := ctx.ShouldBindQuery(params)
	if err != nil {
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    models.Failed,
				"message": "表单验证错误",
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": errs.Error(),
		})
		return
	}

	res := r.logic.OnQueryChildCommentList(params)
	ctx.JSON(http.StatusOK, res)
}

func (r CommentHandler) HandleDelete(ctx *gin.Context) {}

func (r CommentHandler) HandleLike(ctx *gin.Context) {
	form, ok := validate(ctx, models.CommentLikeForm{})
	if !ok {
		return
	}

	res := r.logic.OnCommentLike(form)
	ctx.JSON(http.StatusOK, res)
}
