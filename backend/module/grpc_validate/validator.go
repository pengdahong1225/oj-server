package grpc_validate

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The validate interface starting with protoc-gen-validate v0.6.0.
// See https://github.com/envoyproxy/protoc-gen-validate/pull/455.
type validator interface {
	Validate(all bool) error
}

// The validate interface prior to protoc-gen-validate v0.6.0.
type validatorLegacy interface {
	Validate() error
}

func validate(req interface{}) error {
	switch v := req.(type) {
	case validatorLegacy:
		if err := v.Validate(); err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}
	case validator:
		if err := v.Validate(false); err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}
	}
	return nil
}

// UnaryServerInterceptor 接收到数据后需要检查参数是否合法
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fmt.Println("参数验证...")
		if err := validate(req); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

// UnaryClientInterceptor 发送数据前需要检查参数是否合法
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		fmt.Println("参数验证...")
		if err := validate(req); err != nil {
			return err
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor that validates incoming messages.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapper := &serverWrapper{stream}
		return handler(srv, wrapper)
	}
}

type serverWrapper struct {
	grpc.ServerStream
}

// RecvMsg 接收消息后， 应该先检查一下参数是否合法
func (s *serverWrapper) RecvMsg(m interface{}) error {
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}

	if err := validate(m); err != nil {
		return err
	}

	return nil
}

// SendMsg 发送前， 应该先检查一下参数是否合法
func (s *serverWrapper) SendMsg(m interface{}) error {
	if err := validate(m); err != nil {
		return err
	}
	if err := s.ServerStream.SendMsg(m); err != nil {
		return err
	}
	return nil
}

// clientWrapper  用于包装 grpc.ClientStream 结构体并拦截其对应的方法。
type clientWrapper struct {
	grpc.ClientStream
}

func newWrappedClientStream(c grpc.ClientStream) grpc.ClientStream {
	return &clientWrapper{c}
}

func (c *clientWrapper) RecvMsg(m interface{}) error {
	if err := c.ClientStream.RecvMsg(m); err != nil {
		return err
	}
	if err := validate(m); err != nil {
		return err
	}

	return nil
}

func (c *clientWrapper) SendMsg(m interface{}) error {
	if err := validate(m); err != nil {
		return err
	}
	if err := c.ClientStream.SendMsg(m); err != nil {
		return err
	}
	return nil
}

// ClientStreamInterceptor
func ClientStreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
		streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		s, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, err
		}
		return newWrappedClientStream(s), nil
	}
}
