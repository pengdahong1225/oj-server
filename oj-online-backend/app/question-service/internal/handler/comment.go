package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/logic"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/models"
)

type CommentHandler struct {
}

func (r CommentHandler) Insert(ctx *gin.Context) {
	form, ok := validate(ctx, models.AddCommentForm{})
	if !ok {
		return
	}
	res := logic.CommentLogic{}.OnAddComment(form)
	ctx.JSON(200, res)
}
func (r CommentHandler) Query(ctx *gin.Context) {
	form, ok := validate(ctx, models.QueryCommentForm{})
	if !ok {
		return
	}
	res := logic.CommentLogic{}.OnQueryComment(form)
	ctx.JSON(200, res)
}

func (r CommentHandler) Delete(ctx *gin.Context) {}

func (r CommentHandler) Like(ctx *gin.Context) {}
