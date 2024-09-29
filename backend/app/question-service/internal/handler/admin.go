package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/logic"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/middlewares"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/models"
	"net/http"
	"strconv"
)

type AdminHandler struct{}

func (receiver AdminHandler) UpdateQuestion(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middlewares.UserClaims)
	form, ret := validate(ctx, models.AddProblemForm{})
	if !ret {
		return
	}

	res := logic.AdminLogic{}.HandleUpdateQuestion(claims.Uid, form)
	ctx.JSON(res.Code, res)
}

func (receiver AdminHandler) DeleteQuestion(ctx *gin.Context) {
	p := ctx.GetString("problem_id")
	if p == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
		})
		return
	}
	problemID, _ := strconv.ParseInt(p, 10, 64)
	res := logic.AdminLogic{}.HandleDelQuestion(problemID)
	ctx.JSON(res.Code, res)
}
