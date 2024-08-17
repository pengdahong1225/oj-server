package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/logic"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/middlewares"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/models"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/svc/captcha"
	"net/http"
	"regexp"
	"strconv"
)

type UserHandler struct {
}

func (receiver UserHandler) Login(ctx *gin.Context) {
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

	res := logic.UserLogic{}.HandleLogin(form)
	ctx.JSON(res.Code, res)
}

func (receiver UserHandler) GetUserProfile(ctx *gin.Context) {
	uid, ok := ctx.GetQuery("uid")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
		})
		return
	}

	uidInt, _ := strconv.ParseInt(uid, 10, 64)
	res := logic.UserLogic{}.HandleGetUserProfile(uidInt)
	ctx.JSON(res.Code, res)
}

func (receiver UserHandler) GetRankList(ctx *gin.Context) {
	res := logic.UserLogic{}.HandleGetRankList()
	ctx.JSON(res.Code, res)
}

func (receiver UserHandler) GetSubmitRecord(ctx *gin.Context) {
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
	res := logic.UserLogic{}.HandleGetSubmitRecord(claims.Uid, stamp)
	ctx.JSON(res.Code, res)
}

func (receiver UserHandler) GetUserSolvedList(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middlewares.UserClaims)
	res := logic.UserLogic{}.HandleGetUserSolvedList(claims.Uid)
	ctx.JSON(res.Code, res)
}
