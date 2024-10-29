package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/logic"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/models"
	"net/http"
)

type CommentHandler struct {
}

func (r CommentHandler) HandleAdd(ctx *gin.Context) {
	form, ok := validate(ctx, models.AddCommentForm{})
	if !ok {
		return
	}
	res := logic.CommentLogic{}.OnAddComment(form)
	ctx.JSON(http.StatusOK, res)
}
func (r CommentHandler) HandleGet(ctx *gin.Context) {
	form, ok := validate(ctx, models.QueryCommentForm{})
	if !ok {
		return
	}
	res := logic.CommentLogic{}.OnQueryComment(form)
	ctx.JSON(http.StatusOK, res)
}

func (r CommentHandler) HandleDelete(ctx *gin.Context) {}

func (r CommentHandler) HandleLike(ctx *gin.Context) {}
