package global

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Registry struct {
	client *consulapi.Client
}

func NewRegistry() (*Registry, error) {
	// 配置中心地址
	dsn := fmt.Sprintf("%s:%d", ConfigInstance.Registry_.Host, ConfigInstance.Registry_.Port)
	cfg := consulapi.DefaultConfig()
	cfg.Address = dsn
	// client
	c, err := consulapi.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Registry{client: c}, nil
}

// RegisterService 注册service节点
func (receiver *Registry) RegisterService(serviceName string, ip string, port int) error {
	srv := &consulapi.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", serviceName, ip, port), // 服务唯一ID
		Name:    serviceName,                                    // 服务名称
		Tags:    []string{serviceName},
		Port:    port,
		Address: ip,
		Check: &consulapi.AgentServiceCheck{
			CheckID:                        fmt.Sprintf("%s-%s-%d", serviceName, ip, port),
			HTTP:                           fmt.Sprintf("http://%s:%d%s", ip, port, "/health"),
			Timeout:                        "10s",
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}
	return receiver.client.Agent().ServiceRegister(srv)
}

// NewDBConnection db服务连接
func NewDBConnection() (*grpc.ClientConn, error) {
	registry, err := NewRegistry()
	if err != nil {
		return nil, err
	}
	// registry.client.Health().Service返回的是对应服务的节点列表
	services, _, err := registry.client.Health().Service("db-service", "db-service", true, nil)
	if err != nil {
		return nil, err
	}
	// 这里可以添加简单的负载均衡，访问压力均摊给集群中的每个服务
	dsn := fmt.Sprintf("%s:%d", services[0].Service.Address, services[0].Service.Port)
	return grpc.Dial(dsn, grpc.WithTransportCredentials(insecure.NewCredentials())) // 不安全连接
}
