package registry

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

type ServiceInfo struct {
	Name    string
	NodeId  int
	Address string
	Port    int
	Scheme  string // http/grpc
}

func (s ServiceInfo) ServiceID() string {
	return fmt.Sprintf("%s:%d", s.Name, s.NodeId)
}

type Registry interface {
	// 注册服务，返回注销函数（用于优雅关闭）
	Register(ctx context.Context, inst ServiceInfo) (deregister func() error, err error)

	// 主动注销
	Deregister(ctx context.Context, serviceID string) error

	// 查找服务的健康实例
	Discover(ctx context.Context, serviceName string) ([]ServiceInfo, error)

	// 查找并选择一个实例
	Resolve(ctx context.Context, serviceName string) (ServiceInfo, error)

	// 创建一个到服务的连接
	Dial(ctx context.Context, name string) (*grpc.ClientConn, error)

	// 关闭客户端
	Close() error
}
