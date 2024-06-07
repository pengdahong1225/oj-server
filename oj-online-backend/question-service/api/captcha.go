package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"question-service/api/internal"
	"regexp"
)

// 图像验证码
func GetImageCode(ctx *gin.Context) {}

// 短信验证码
func GetSmsCode(ctx *gin.Context) {
	// 手机号校验
	if mobile, ok := ctx.GetQuery("mobile"); !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
		})
		ctx.Abort()
		return
	} else {
		ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, mobile)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "参数错误",
			})
			ctx.Abort()
		}
		if err := internal.SendSmsCode(mobile); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": "发送成功[测试环境的验证码均为123456]",
			})
		}
	}
}
