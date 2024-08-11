package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/api/handler"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/models"
)

type Comment struct {
}

func (r Comment) Add(ctx *gin.Context) {
	form, ok := validate(ctx, models.CommentForm{})
	if !ok {
		return
	}
	res := handler.CommentHandler{}.HandleAddComment(form)
	ctx.JSON(200, res)
}

func (r Comment) Delete(ctx *gin.Context) {}

func (r Comment) Like(ctx *gin.Context) {}

func (r Comment) Get(ctx *gin.Context) {}
