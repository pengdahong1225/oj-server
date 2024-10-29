package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/logic"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/middlewares"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/models"
	"net/http"
	"regexp"
	"strconv"
)

type User struct {
	logic logic.User
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
	res := r.logic.OnUserRegister(form)
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

	res := r.logic.OnUserLogin(form)
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
	res := r.logic.GetUserProfile(uidInt)
	ctx.JSON(http.StatusOK, res)
}

func (r User) HandleRankList(ctx *gin.Context) {
	res := r.logic.GetRankList()
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
	res := r.logic.GetSubmitRecord(claims.Uid, stamp)
	ctx.JSON(http.StatusOK, res)
}

func (r User) HandleSolvedList(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middlewares.UserClaims)
	res := r.logic.GetUserSolvedList(claims.Uid)
	ctx.JSON(http.StatusOK, res)
}

func (r User) HandleUPSS(ctx *gin.Context) {
	p := ctx.Query("params")
	if p == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "参数错误",
		})
		ctx.Abort()
		return
	}

	params := &models.UPSSParams{}
	err := json.Unmarshal([]byte(p), params)
	if err != nil || len(params.ProblemIds) <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": err.Error(),
		})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, r.logic.QueryUserSolvedListByProblemList(params))
}
