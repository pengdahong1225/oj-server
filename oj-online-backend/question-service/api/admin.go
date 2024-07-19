package api

import (
	"github.com/gin-gonic/gin"
	"question-service/api/internal"
	"question-service/middlewares"
	"question-service/models"
)

func AddQuestion(ctx *gin.Context) {
	// 表单校验
	claims := ctx.MustGet("claims").(*middlewares.UserClaims)
	form, ret := validate(ctx, models.AddProblemForm{})
	if !ret {
		return
	}

	res := internal.AdminHandler{}.HandleAddQuestion(claims.Uid, form)
	ctx.JSON(res.Code, res)
}
func DeleteQuestion(ctx *gin.Context) {}
func UpdateQuestion(ctx *gin.Context) {

}
