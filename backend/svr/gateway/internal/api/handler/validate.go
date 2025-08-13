package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"oj-server/proto/pb"
	"oj-server/svr/gateway/internal/api/define"
)

// 表单类型集
type formTyper interface {
	define.RegisterForm | define.LoginFrom | define.GetSmsCodeForm | define.SubmitForm | define.UpdateProblemForm |
		define.AddCommentForm | define.CommentLikeForm | define.LoginWithSmsForm | define.ResetPasswordForm |
		define.CreateProblemForm
}

// validate 通用表单验证
func validate[T formTyper](ctx *gin.Context, form T) (*T, bool) {
	resp := &define.Response{
		ErrCode: pb.Error_EN_Success,
		Message: "",
	}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		errs, ok := err.(validator.ValidationErrors)
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
