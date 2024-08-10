package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/api/handler"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/middlewares"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/models"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/services/captcha"
	"net/http"
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

	res := handler.UserHandler{}.HandleLogin(form)
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
	res := handler.UserHandler{}.HandleGetUserProfile(uidInt)
	ctx.JSON(res.Code, res)
}

func GetRankList(ctx *gin.Context) {
	res := handler.UserHandler{}.HandleGetRankList()
	ctx.JSON(res.Code, res)
}

func GetSubmitRecord(ctx *gin.Context) {
	t, ok := ctx.GetQuery("stamp")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
		})
		return
	}
	claims := ctx.MustGet("claims").(*middlewares.UserClaims)
	stamp, _ := strconv.ParseInt(t, 10, 64)
	res := handler.UserHandler{}.HandleGetSubmitRecord(claims.Uid, stamp)
	ctx.JSON(res.Code, res)
}

func GetUserSolvedList(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middlewares.UserClaims)
	res := handler.UserHandler{}.HandleGetUserSolvedList(claims.Uid)
	ctx.JSON(res.Code, res)
}
