package middlewares

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()

	// 记录请求信息
	logrus.Infof("gRPC 请求开始 | 方法: %s | 请求参数: %+v", info.FullMethod, req)

	// 调用后续处理链
	resp, err = handler(ctx, req)

	// 记录响应信息
	logrus.Infof("gRPC 请求结束 | 方法: %s | 耗时: %v | 错误: %v",
		info.FullMethod, time.Since(start), err)

	return resp, err
}
