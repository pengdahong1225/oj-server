package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/logic"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/middlewares"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/models"
	"net/http"
	"regexp"
	"strconv"
)

type User struct {
}

func (r User) HandleRegister(ctx *gin.Context) {
	// 表单验证
	form, ret := validate(ctx, models.RegisterForm{})
	if !ret {
		return
	}

	// 手机号校验
	ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "手机号格式错误",
		})
		return
	}
	// 密码校验
	if form.PassWord != form.RePassWord {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "两次密码输入不匹配",
		})
		return
	}
	res := logic.User{}.OnUserRegister(form)
	ctx.JSON(http.StatusOK, res)
}
func (r User) HandleLogin(ctx *gin.Context) {
	// 表单验证
	form, ret := validate(ctx, models.LoginFrom{})
	if !ret {
		return
	}

	// 手机号校验
	ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, form.Mobile)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "手机号格式错误",
		})
		return
	}

	res := logic.User{}.HandleLogin(form)
	ctx.JSON(http.StatusOK, res)
}
func (r User) HandleResetPassword(ctx *gin.Context) {

}

func (r User) HandleUserProfile(ctx *gin.Context) {
	uid, ok := ctx.GetQuery("uid")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "参数错误",
		})
		return
	}

	uidInt, _ := strconv.ParseInt(uid, 10, 64)
	res := logic.User{}.HandleGetUserProfile(uidInt)
	ctx.JSON(http.StatusOK, res)
}

func (r User) HandleRankList(ctx *gin.Context) {
	res := logic.User{}.HandleGetRankList()
	ctx.JSON(http.StatusOK, res)
}

func (r User) HandleSubmitRecord(ctx *gin.Context) {
	t, ok := ctx.GetQuery("stamp")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "参数错误",
		})
		return
	}
	claims := ctx.MustGet("claims").(*middlewares.UserClaims)
	stamp, _ := strconv.ParseInt(t, 10, 64)
	res := logic.User{}.HandleGetSubmitRecord(claims.Uid, stamp)
	ctx.JSON(http.StatusOK, res)
}

func (r User) HandleSolvedList(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middlewares.UserClaims)
	res := logic.User{}.HandleGetUserSolvedList(claims.Uid)
	ctx.JSON(http.StatusOK, res)
}
