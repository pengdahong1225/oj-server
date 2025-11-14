package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 统一返回格式
type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func ResponseWithJson(ctx *gin.Context, code int, message string, data any) {
	ctx.JSON(code, APIResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// 带文件流成功响应
func RespondWithBytesAsFile(ctx *gin.Context, data []byte, contentType string, downloadName string) {
	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename=%s`, downloadName))
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(data)))

	ctx.Data(http.StatusOK, contentType, data)
}
