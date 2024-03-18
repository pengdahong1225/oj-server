package global

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
)

type Registry struct {
	client *consulapi.Client
}

func NewRegistry() (*Registry, error) {
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

func (receiver *Registry) RegisterService(serviceName string, ip string, port int) error {
	srv := &consulapi.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", serviceName, ip, port), // 服务唯一ID
		Name:    serviceName,                                    // 服务名称
		Tags:    []string{serviceName},
		Port:    port,
		Address: ip,
		// 注意：TCP健康检查无需提供CheckFunc，因为它是由Consul Agent发起的TCP连接尝试
		Check: &consulapi.AgentServiceCheck{
			CheckID:                        fmt.Sprintf("%s-%s-%d", serviceName, ip, port),
			TCP:                            fmt.Sprintf("%s:%d", ip, port),
			Timeout:                        "10s",
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}
	return receiver.client.Agent().ServiceRegister(srv)
}

// NewQuestionConnection Question服务连接
func QuestionDsn() (string, error) {
	registry, err := NewRegistry()
	if err != nil {
		return "", err
	}
	// registry.client.Health().Service返回的是对应服务的节点列表
	services, _, err := registry.client.Health().Service("question-service", "question-service", true, nil)
	if err != nil {
		return "", err
	}
	// 这里可以添加简单的负载均衡，访问压力均摊给集群中的每个服务
	dsn := fmt.Sprintf("http://%s:%d", services[0].Service.Address, services[0].Service.Port)
	return dsn, nil
}
