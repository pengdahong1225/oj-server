package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"github.com/pengdahong1225/oj-server/backend/app/gateway/internal/define"
)

// 表单类型集
type formTyper interface {
	define.RegisterForm | define.LoginFrom | define.GetSmsCodeForm | define.SubmitForm | define.UpdateProblemForm |
	define.AddCommentForm | define.NoticeForm | define.CommentLikeForm
}

// validate 通用表单验证
func validate[T formTyper](ctx *gin.Context, form T) (*T, bool) {
	if err := ctx.ShouldBindJSON(&form); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    define.Failed,
				"message": "表单验证错误",
			})
			return nil, false
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    define.Failed,
			"message": errs.Error(),
		})
		return nil, false
	}
	return &form, true
}
