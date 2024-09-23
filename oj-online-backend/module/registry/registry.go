package registry

import (
	"fmt"
	"github.com/pengdahong1225/Oj-Online-Server/module/settings"

	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
func (receiver *Registry) RegisterServiceWithGrpc(serviceName string, ip string, port int) error {
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
func (receiver *Registry) RegisterServiceWithHttp(serviceName string, ip string, port int) error {
	srv := &consulapi.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", serviceName, ip, port), // 服务唯一ID
		Name:    serviceName,                                    // 服务名称
		Tags:    []string{serviceName},
		Port:    port,
		Address: ip,
		Check: &consulapi.AgentServiceCheck{
			CheckID:                        fmt.Sprintf("%s-%s-%d", serviceName, ip, port),
			HTTP:                           fmt.Sprintf("%s:%d/%s", ip, port, "health"),
			Timeout:                        "10s",
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}
	return receiver.client.Agent().ServiceRegister(srv)
}

// NewDBConnection db服务连接
func NewDBConnection(cfg *settings.RegistryConfig) (*grpc.ClientConn, error) {
	register, err := NewRegistry(cfg)
	if err != nil {
		return nil, err
	}
	// registry.client.Health().Service返回的是对应服务的节点列表
	services, _, err := register.client.Health().Service("db-service", "db-service", true, nil)
	if err != nil {
		return nil, err
	}
	// 如果是集群的话，这里可以添加简单的负载均衡，访问压力均摊给集群中的每个服务
	dsn := fmt.Sprintf("%s:%d", services[0].Service.Address, services[0].Service.Port)
	return grpc.NewClient(dsn, grpc.WithTransportCredentials(insecure.NewCredentials())) // 不安全连接
}

// NewJudgeConnection judge服务连接
func GetJudgeServerDsn(cfg *settings.RegistryConfig) (string, error) {
	register, err := NewRegistry(cfg)
	if err != nil {
		return "", err
	}
	// registry.client.Health().Service返回的是对应服务的节点列表
	services, _, err := register.client.Health().Service("judge-service", "judge-service", true, nil)
	if err != nil {
		return "", err
	}
	// 如果是集群的话，这里可以添加简单的负载均衡，访问压力均摊给集群中的每个服务
	dsn := fmt.Sprintf("%s:%d", services[0].Service.Address, services[0].Service.Port)
	return dsn, nil
}
