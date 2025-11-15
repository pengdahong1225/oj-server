package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
	"net/http"
	"oj-server/pkg/proto/pb"
	"strings"
)

// 统一返回格式
type APIResponse struct {
	Code    pb.Error `json:"code"`
	Message string   `json:"message"`
	Data    any      `json:"data"`
}

// --- 通用成功 ---
func ResponseOK(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    pb.Error_EN_Success,
		Message: "success",
		Data:    data,
	})
}

// --- 通用失败（业务错误） ---
func ResponseError(ctx *gin.Context, code pb.Error, msg string) {
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    code,
		Message: msg,
		Data:    nil,
	})
}

// --- gRPC 错误转 REST 错误（安全） ---
func ResponseWithGrpcError(ctx *gin.Context, err error) {
	st, ok := status.FromError(err)
	if ok {
		// 只返回业务 message，不返回 grpc 内部错误
		msg := sanitizeGrpcMessage(st.Message())

		ctx.JSON(http.StatusOK, APIResponse{
			Code:    pb.Error(st.Code()),
			Message: msg,
			Data:    nil,
		})
		return
	}

	// 未知错误
	ResponseError(ctx, pb.Error_EN_Failed, "未知错误")
}

// --- 过滤掉 gRPC 内部调试信息 ---
func sanitizeGrpcMessage(msg string) string {
	if strings.HasPrefix(msg, "rpc error:") {
		return "服务内部错误"
	}
	return msg
}

// 带文件流成功响应
func RespondWithBytesAsFile(ctx *gin.Context, data []byte, contentType string, downloadName string) {
	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename=%s`, downloadName))
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(data)))

	ctx.Data(http.StatusOK, contentType, data)
}

// HTTP 错误响应
func ResponseBadRequest(ctx *gin.Context, msg string) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"error":   "bad_request",
		"message": msg,
	})
}

func ResponseUnauthorized(ctx *gin.Context, msg string) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error":   "unauthorized",
		"message": msg,
	})
}

func ResponseForbidden(ctx *gin.Context, msg string) {
	ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"error":   "forbidden",
		"message": msg,
	})
}

func ResponseInternalServerError(ctx *gin.Context, msg string) {
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"error":   "internal_server_error",
		"message": msg,
	})
}
