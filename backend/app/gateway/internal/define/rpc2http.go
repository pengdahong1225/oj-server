package define

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

// error_mapping.go
var grpcToHTTP = map[codes.Code]int{
	codes.OK:                 http.StatusOK,
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.Unauthenticated:    http.StatusUnauthorized,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
	codes.FailedPrecondition: http.StatusPreconditionFailed,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,
	codes.Internal:           http.StatusInternalServerError,
}

func ConvertGrpcCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	// 提取gRPC状态码
	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError
	}

	// 转换HTTP状态码
	httpCode := grpcToHTTP[st.Code()]
	if httpCode == 0 {
		httpCode = http.StatusInternalServerError
	}

	return httpCode
}
