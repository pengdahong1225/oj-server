package registry

import (
	"fmt"
	"github.com/pengdahong1225/Oj-Online-Server/module/settings"
	"github.com/pengdahong1225/Oj-Online-Server/module/signal"
	"github.com/pengdahong1225/Oj-Online-Server/module/utils"

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
	nodeID, err := utils.GenerateUUID()
	if err != nil {
		return err
	}
	srv := &consulapi.AgentServiceRegistration{
		ID:      nodeID,      // 服务唯一ID
		Name:    serviceName, // 服务名称
		Tags:    []string{serviceName},
		Port:    port,
		Address: ip,
		Check: &consulapi.AgentServiceCheck{
			CheckID:                        nodeID,
			GRPC:                           fmt.Sprintf("%s:%d", ip, port),
			Interval:                       "5s",  // 每5秒检测一次
			Timeout:                        "5s",  // 5秒超时
			DeregisterCriticalServiceAfter: "10s", // 超时10秒注销节点
		},
	}
	err = receiver.client.Agent().ServiceRegister(srv)
	if err != nil {
		return err
	}
	go signal.SignalListen(func() {
		receiver.client.Agent().ServiceDeregister(nodeID)
	})
	return nil
}
func (receiver *Registry) RegisterServiceWithHttp(serviceName string, ip string, port int) error {
	nodeID, err := utils.GenerateUUID()
	if err != nil {
		return err
	}
	srv := &consulapi.AgentServiceRegistration{
		ID:      nodeID,      // 服务唯一ID
		Name:    serviceName, // 服务名称
		Tags:    []string{serviceName},
		Port:    port,
		Address: ip,
		Check: &consulapi.AgentServiceCheck{
			CheckID:                        nodeID,
			HTTP:                           fmt.Sprintf("http://%s:%d/%s", ip, port, "health"),
			Interval:                       "5s",  // 每5秒检测一次
			Timeout:                        "5s",  // 5秒超时
			DeregisterCriticalServiceAfter: "10s", // 超时10秒注销节点
		},
	}
	err = receiver.client.Agent().ServiceRegister(srv)
	if err != nil {
		return err
	}
	go signal.SignalListen(func() {
		receiver.client.Agent().ServiceDeregister(nodeID)
	})
	return nil
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
