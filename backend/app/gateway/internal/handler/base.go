package handler

import "net/http"

// HTTPHandler 定义处理HTTP请求的模板
type HTTPHandler interface {
	Authenticate(r *http.Request) error
	Validate(r *http.Request) error
	Process(r *http.Request) (interface{}, error)
	FormatResponse(data interface{}) []byte
}

// BaseHTTPHandler 基础实现
type BaseHTTPHandler struct {
	HTTPHandler
}

// HandleRequest 模板方法
func (h *BaseHTTPHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if err := h.Authenticate(r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := h.Validate(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := h.Process(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(h.FormatResponse(data))
}
