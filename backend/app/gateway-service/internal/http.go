package internal

import (
	"context"
	"encoding/json"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"net/http"
)

// HttpResponse
// http 响应结构体
type HttpResponse struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func GatewayResponseModifier(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	// 返回的数据，在外层同一封装了数据结构HTTPResponse，对一些历史项目兼容有很棒的效果
	newResp := &HttpResponse{
		Data: resp,
	}
	pbData, _ := json.Marshal(newResp)
	_, _ = w.Write(pbData)
	return nil
}

// error尽可能用gRPC标准的错误Status表示
// gRPC的标准错误，对错误码code有一套定义（参考google.golang.org/grpc/codes），类似于HTTP的状态码
// 错误码code要尽量少，过多没有意义
//	标准错误码尽可能复用，如资源找不到、权限不足等
//	业务错误码可以独立，一般一个系统定义1个即可

func WithErrorHandler(fn ErrorHandlerFunc) ServeMuxOption {

}
func GatewayErrModifier(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	// 提取error
	s, ok := status.FromError(err)
	// 非标准错误
	if !ok {
		runtime.DefaultHTTPErrorHandler(ctx, mux, m, w, r, err)
		return
	}

	// 对各类错误增加定制的逻辑
	switch s.Code() {
	case codes.Unauthenticated:
		// 示例：认证失败，可以加入重定向的逻辑
	default:
		runtime.DefaultHTTPErrorHandler(ctx, mux, m, w, r, err)
	}

	return
}
