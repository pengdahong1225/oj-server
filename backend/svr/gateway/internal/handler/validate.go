package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"oj-server/svr/gateway/internal/model"
)

type formTyper interface {
	model.RegisterForm | model.LoginFrom | model.GetSmsCodeForm |
		model.CommentLikeForm | model.LoginWithSmsForm | model.ResetPasswordForm |
		model.DeleteCommentForm
}

func validateWithForm[T formTyper](ctx *gin.Context, form T) (*T, bool) {
	if err := ctx.ShouldBind(&form); err != nil {
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseBadRequest(ctx, "表单验证-未知错误")
			return nil, false
		}
		ResponseBadRequest(ctx, errs.Error())
		return nil, false
	}
	return &form, true
}

type jsonTyper interface {
	model.SubmitForm | model.CreateCommentForm | model.CreateProblemForm | model.UpdateProblemForm
}

func validateWithJson[T jsonTyper](ctx *gin.Context, form T) (*T, bool) {
	if err := ctx.ShouldBindJSON(&form); err != nil {
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseBadRequest(ctx, "表单验证-未知错误")
			return nil, false
		}
		ResponseBadRequest(ctx, errs.Error())
		return nil, false
	}
	return &form, true
}
