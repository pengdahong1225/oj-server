package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"question-service/api/internal"
	"question-service/models"
	"question-service/services/captcha"
	"regexp"
	"strconv"
)

func Login(ctx *gin.Context) {
	// 表单验证
	form, ret := validate(ctx, models.LoginFrom{})
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

	// 短信验证码校验
	if !captcha.VerifySmsCode(form.Mobile, form.SmsCode) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "验证码错误",
		})
		return
	}

	res := internal.ProcessForLogin(form)
	ctx.JSON(res.Code, res)
}

func GetUserDetail(ctx *gin.Context) {
	// 查询参数
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
			return
		}
		mobileInt, _ := strconv.ParseInt(mobile, 10, 64)
		res := internal.GetUserDetail(mobileInt)
		ctx.JSON(res.Code, res)
	}
}

func GetRankList(ctx *gin.Context) {
	res := internal.GetRankList()
	ctx.JSON(res.Code, res)
}

func GetSubmitRecord(ctx *gin.Context) {

}
