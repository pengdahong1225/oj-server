package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"question-service/api/internal"
	"regexp"
	"strconv"
)

func UserRegister(ctx *gin.Context) {
	// 表单验证
	if form, ok := formValidateForRegistry(ctx); ok {
		internal.RegistryHandler(ctx, form)
	}
}

func UserLogin(ctx *gin.Context) {
	// 表单验证
	if form, ok := formValidateForLogin(ctx); ok {
		internal.LoginHandler(ctx, form)
	}
}

func GetUserDetail(ctx *gin.Context) {
	// 查询参数
	if phone, ok := ctx.GetQuery("phone"); !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		ctx.Abort()
		return
	} else {
		ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, phone)
		if !ok {
			ctx.Abort()
			return
		}
		phoneInt, _ := strconv.ParseInt(phone, 10, 64)
		internal.GetUserDetail(ctx, phoneInt)
	}
}

func GetRankList(ctx *gin.Context) {
	internal.GetRankList(ctx)
}

func GetSubmitRecord(ctx *gin.Context) {
	if id, ok := ctx.GetQuery("userId"); !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		ctx.Abort()
		return
	} else {
		userId, _ := strconv.ParseInt(id, 10, 64)
		internal.GetSubmitRecord(ctx, userId)
	}
}

func SendSmsCode(ctx *gin.Context) {
	// 手机号校验
	if phone, ok := ctx.GetQuery("phone"); !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		ctx.Abort()
		return
	} else {
		ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, phone)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "参数错误",
			})
			ctx.Abort()
		}
		if err := internal.SendSmsCode(phone); err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"msg": "发送成功",
			})
		}
	}
}
