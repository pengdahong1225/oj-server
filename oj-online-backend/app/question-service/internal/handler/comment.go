package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/logic"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/models"
)

type CommentHandler struct {
}

func (r CommentHandler) Add(ctx *gin.Context) {
	form, ok := validate(ctx, models.CommentForm{})
	if !ok {
		return
	}
	res := logic.CommentLogic{}.HandleAddComment(form)
	ctx.JSON(200, res)
}

func (r CommentHandler) Delete(ctx *gin.Context) {}

func (r CommentHandler) Like(ctx *gin.Context) {}

func (r CommentHandler) Get(ctx *gin.Context) {}
