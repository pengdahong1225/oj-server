package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
	"question-service/models"
	"regexp"
)

// 表单类型集
type formTyper interface {
	models.LoginFrom | models.RegistryForm | models.QuestionForm
}

// validate 表单验证
func validate[T formTyper](ctx *gin.Context, form T) (*T, bool) {
	if err := ctx.ShouldBindJSON(&form); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ctx.JSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"message": "表单验证错误",
			})
			return nil, false
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": errs.Error(),
		})
		return nil, false
	}
	return &form, true
}

// validateForLogin 登录表单验证
// gin无法校验手机号格式，需要自定制
// 修改gin中的Validator引擎属性，注册新的校验函数
func validateForLogin(ctx *gin.Context) (*models.LoginFrom, bool) {
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := validate.RegisterValidation("mobile", func(fl validator.FieldLevel) bool {
			mobile := fl.Field().String()
			ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, mobile)
			if !ok {
				return false
			}
			return true
		}); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
			return nil, false
		}
	}
	passwordLoginForm := models.LoginFrom{}
	return validate(ctx, passwordLoginForm)
}
