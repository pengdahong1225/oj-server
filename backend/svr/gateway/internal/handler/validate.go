package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"oj-server/module/proto/pb"
	"oj-server/svr/gateway/internal/model"
)

type formTyper interface {
	model.RegisterForm | model.LoginFrom | model.GetSmsCodeForm | model.UpdateProblemForm |
		model.AddCommentForm | model.CommentLikeForm | model.LoginWithSmsForm | model.ResetPasswordForm |
		model.CreateProblemForm
}

func validateWithForm[T formTyper](ctx *gin.Context, form T) (*T, bool) {
	resp := &model.Response{
		ErrCode: pb.Error_EN_Success,
		Message: "",
	}
	if err := ctx.ShouldBind(&form); err != nil {
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			resp.ErrCode = pb.Error_EN_Failed
			resp.Message = "表单验证-未知错误"
			ctx.JSON(http.StatusBadRequest, resp)
			return nil, false
		}
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = errs.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return nil, false
	}
	return &form, true
}

type jsonTyper interface {
	model.SubmitForm
}

func validateWithJson[T jsonTyper](ctx *gin.Context, form T) (*T, bool) {
	resp := &model.Response{
		ErrCode: pb.Error_EN_Success,
		Message: "",
	}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			resp.ErrCode = pb.Error_EN_Failed
			resp.Message = "表单验证-未知错误"
			ctx.JSON(http.StatusBadRequest, resp)
			return nil, false
		}
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = errs.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return nil, false
	}
	return &form, true
}
