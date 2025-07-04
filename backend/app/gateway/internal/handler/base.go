package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HTTPHandler 定义处理HTTP请求的模板
type HTTPHandler interface {
	Authenticate(ctx *gin.Context) error           // 鉴权
	Validate(ctx *gin.Context) error               // 参数校验
	Process(ctx *gin.Context) (interface{}, error) // 处理请求
	FormatResponse(data any) []byte                // 格式化响应
}

// BaseHTTPHandler 基础实现
type BaseHTTPHandler struct {
	HTTPHandler
}

// HandleRequest 模板方法
func (h *BaseHTTPHandler) HandleRequest(ctx *gin.Context) {
	err := h.Authenticate(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = h.Validate(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.Process(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Data(http.StatusOK, "application/json", h.FormatResponse(data))
}
