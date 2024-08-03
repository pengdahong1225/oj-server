package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"question-service/internal/api/handler"
	"question-service/internal/middlewares"
	"question-service/models"
	"strconv"
)

func UpdateQuestion(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middlewares.UserClaims)
	form, ret := validate(ctx, models.AddProblemForm{})
	if !ret {
		return
	}

	res := handler.AdminHandler{}.HandleUpdateQuestion(claims.Uid, form)
	ctx.JSON(res.Code, res)
}

func DeleteQuestion(ctx *gin.Context) {
	p := ctx.GetString("problem_id")
	if p == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
		})
		return
	}
	problemID, _ := strconv.ParseInt(p, 10, 64)
	res := handler.AdminHandler{}.HandleDelQuestion(problemID)
	ctx.JSON(res.Code, res)
}
