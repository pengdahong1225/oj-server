package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"question-service/api/internal"
	"strconv"
)

func GetUserList(ctx *gin.Context) {
	cursor, _ := strconv.Atoi(ctx.DefaultQuery("cursor", "0"))
	if cursor < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		ctx.Abort()
		return
	}
	internal.GetUserList(ctx, int32(cursor))
}

func AddQuestion(ctx *gin.Context) {

}
func DeleteQuestion(ctx *gin.Context) {

}
func UpdateQuestion(ctx *gin.Context) {

}
