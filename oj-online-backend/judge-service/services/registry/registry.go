package registry

import (
	"db-service/settings"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
)

type Registry struct {
	client *consulapi.Client
}

func NewRegistry(cfg *settings.RegistryConfig) (*Registry, error) {
	// 配置中心地址
	dsn := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	consulConf := consulapi.DefaultConfig()
	consulConf.Address = dsn
	// client
	c, err := consulapi.NewClient(consulConf)
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
