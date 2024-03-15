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
		Check: &consulapi.AgentServiceCheck{
			CheckID:                        fmt.Sprintf("%s-%s-%d", serviceName, ip, port),
			GRPC:                           fmt.Sprintf("%s:%d", ip, port),
			Timeout:                        "10s",
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}
	return receiver.client.Agent().ServiceRegister(srv)
}
