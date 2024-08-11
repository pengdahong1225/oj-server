package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/models"
	"net/http"
)

// 表单类型集
type formTyper interface {
	models.LoginFrom | models.GetSmsCodeForm | models.SubmitForm | models.AddProblemForm
}

// validate 通用表单验证
func validate[T formTyper](ctx *gin.Context, form T) (*T, bool) {
	if err := ctx.ShouldBindJSON(&form); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "表单验证错误",
			})
			return nil, false
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": errs.Error(),
		})
		return nil, false
	}
	return &form, true
}
