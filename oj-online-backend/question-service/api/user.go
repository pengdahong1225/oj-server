package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"question-service/api/internal"
	"question-service/middlewares"
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

func GetUserProfile(ctx *gin.Context) {
	uid, ok := ctx.GetQuery("uid")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
		})
		return
	}

	uidInt, _ := strconv.ParseInt(uid, 10, 64)
	res := internal.GetUserProfileByUid(uidInt)
	ctx.JSON(res.Code, res)
}

func GetRankList(ctx *gin.Context) {
	res := internal.GetRankList()
	ctx.JSON(res.Code, res)
}

func GetSubmitRecord(ctx *gin.Context) {

}

func GetUserSolvedList(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middlewares.UserClaims)
	res := internal.GetUserSolvedList(claims.Uid)
	ctx.JSON(res.Code, res)
}
