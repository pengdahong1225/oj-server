package middlewares

import (
	"errors"
	"github.com/pengdahong1225/oj-server/backend/app/gateway/internal/define"
	"net/http"
)

// middleware.go
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 使用自定义ResponseWriter捕获响应
		rw := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(rw, r)

		// 仅处理错误情况
		if rw.status < 400 {
			return
		}

		// 解析gRPC错误
		httpCode := define.ConvertGrpcCode(
			errors.New(string(rw.body)),
		)

		// 返回标准化错误
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpCode)
	})
}

// 自定义ResponseWriter
type responseWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body = b
	return rw.ResponseWriter.Write(b)
}
