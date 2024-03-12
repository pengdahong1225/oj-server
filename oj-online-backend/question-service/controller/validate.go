package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
	"question-service/models"
	"regexp"
)

func formValidateForRegistry(ctx *gin.Context) (*models.RegistryForm, bool) {
	// 手机号 -- 修改gin框架中的Validator引擎属性，实现自定制
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := validate.RegisterValidation("phone", validatePhone); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return nil, false
		}
	}
	// 密码
	registryForm := models.RegistryForm{}
	return processOnValidate(ctx, registryForm)
}

// 表单验证
func formValidateForLogin(ctx *gin.Context) (*models.LoginFrom, bool) {
	// 手机号 -- 修改gin框架中的Validator引擎属性，实现自定制
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义字段级别校验方法
		if err := validate.RegisterValidation("phone", validatePhone); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return nil, false
		}
	}
	// 密码
	passwordLoginForm := models.LoginFrom{}
	return processOnValidate(ctx, passwordLoginForm)
}

// 自定义验证函数：手机号格式校验 使用正则表达式
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, phone)
	if !ok {
		return false
	}
	return true
}

// 类型集
type formTyper interface {
	models.LoginFrom | models.RegistryForm | models.QuestionForm
}

func processOnValidate[T formTyper](ctx *gin.Context, form T) (*T, bool) {
	if err := ctx.ShouldBindJSON(&form); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ctx.JSON(http.StatusForbidden, gin.H{"msg": "表单验证错误"})
			return nil, false
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": errs.Error(),
		})
		return nil, false
	}
	return &form, true
}
