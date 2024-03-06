package controller

import (
	"github.com/gin-gonic/gin"
	"question-service/logic"
)

func UserRegister(ctx *gin.Context) {
	// 表单验证
	if form, ok := formValidateForRegistry(ctx); ok {
		logic.RegistryHandler(ctx, form)
	}
}

func UserLogin(ctx *gin.Context) {

}

func GetUserDetail(ctx *gin.Context) {

}

func GetRankList(ctx *gin.Context) {

}

func GetSubmitRecord(ctx *gin.Context) {

}

func SendCmsCode(ctx *gin.Context) {

}
