package handler

import (
	"net/http"
	"oj-server/module/captcha"
	"oj-server/proto/pb"
	"oj-server/src/gateway/internal/data"
	"oj-server/src/gateway/internal/define"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 获取图像验证码
func HandleGetImageCode(ctx *gin.Context) {
	resp := define.Response{
		ErrCode: pb.Error_EN_Success,
	}
	id, b64s, err := captcha.GenerateImageCaptcha()
	if err != nil {
		logrus.Errorf("生成图像验证码失败: %v", err)
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务繁忙"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	data := define.ImageCaptchaData{
		CaptchaID: id,
		Captcha:   b64s,
	}
	resp.ErrCode = pb.Error_EN_Success
	resp.Data = data
	resp.Message = "图形验证码生成成功"
	ctx.JSON(http.StatusOK, resp)
}

// 获取短信验证码
func HandleGetSmsCode(ctx *gin.Context) {
	// 表单校验
	form, ret := validate(ctx, define.GetSmsCodeForm{})
	if !ret {
		return
	}

	resp := define.Response{
		ErrCode: pb.Error_EN_Success,
	}
	// 手机号校验
	ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "手机号格式错误"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	// 图形验证码校验
	if !captcha.VerifyImageCaptcha(form.CaptchaID, form.CaptchaValue) {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "图形验证码输入错误"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// 发送验证码
	c, err := captcha.SendSmsCode(form.Mobile)
	if err != nil {
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务繁忙"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	// 缓存验证码
	err = data.SetSmsCaptcha(form.Mobile, c)
	if err != nil {
		logrus.Errorf("缓存验证码失败: %v", err)
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务繁忙"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.ErrCode = pb.Error_EN_Success
	resp.Message = "短信验证码发送成功"
	ctx.JSON(http.StatusOK, resp)
}
