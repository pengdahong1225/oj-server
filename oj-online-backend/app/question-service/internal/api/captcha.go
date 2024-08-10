package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/models"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/services/captcha"
	"net/http"
	"regexp"
)

// GetImageCode 获取图像验证码
func GetImageCode(ctx *gin.Context) {
	id, b64s, err := captcha.GenerateImageCaptcha()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
	} else {
		data := make(map[string]any)
		data["captchaID"] = id
		data["captchaUrl"] = b64s
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "OK",
			"data":    data,
		})
	}
}

// GetSmsCode 获取短信验证码
func GetSmsCode(ctx *gin.Context) {
	// 表单校验
	form, ret := validate(ctx, models.GetSmsCodeForm{})
	if !ret {
		return
	}

	// 手机号校验
	ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "手机号格式错误",
		})
		return
	}
	// 图形验证码校验
	if captcha.VerifyImageCaptcha(form.CaptchaID, form.CaptchaValue) {
		if err := captcha.SendSmsCode(form.Mobile); err != nil {
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
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "图形验证码输入错误",
		})
	}
}
