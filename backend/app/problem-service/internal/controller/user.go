package controller

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
	uid := ctx.GetInt64("uid")
	res := r.logic.GetUserProfile(uid)
	ctx.JSON(http.StatusOK, res)
}

func (r User) HandleRankList(ctx *gin.Context) {
	res := r.logic.GetRankList()
	ctx.JSON(http.StatusOK, res)
}

func (r User) HandleRefreshToken(ctx *gin.Context) {
	token := ctx.Request.Header.Get("token")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    models.Failed,
			"message": "未登录",
		})
		return
	}
	j := middlewares.NewJWT()
	newToken, err := j.RefreshToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    models.Failed,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    models.Success,
		"message": "OK",
		"token":   newToken,
	})
}

// HandleRecordList 历史提交记录
func (r User) HandleRecordList(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("page_size")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "页码参数错误",
		})
		ctx.Abort()
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "页大小参数错误",
		})
		ctx.Abort()
		return
	}

	uid := ctx.GetInt64("uid")
	res := r.logic.GetRecordList(uid, page, pageSize)
	ctx.JSON(http.StatusOK, res)
}

func (r User) HandleRecord(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": "参数错误",
		})
		ctx.Abort()
		return
	}

	res := r.logic.GetRecord(int64(id))
	ctx.JSON(http.StatusOK, res)
}

func (r User) HandleSolvedList(ctx *gin.Context) {
	uid := ctx.GetInt64("uid")
	res := r.logic.GetUserSolvedList(uid)
	ctx.JSON(http.StatusOK, res)
}
