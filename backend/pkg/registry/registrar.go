package registry

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"sync"
)

var (
	MyRegistrar *Registrar
)

type ServiceInfo struct {
	Name    string
	NodeId  int
	Address string
	Port    int
	Scheme  string // http/grpc
}

type Registrar struct {
	dsn    string
	client *consulapi.Client

	servicesMap map[string]*grpc.ClientConn // 服务连接池
	mux         sync.RWMutex
}

func NewRegistrar(dsn string) (*Registrar, error) {
	consulConf := consulapi.DefaultConfig()
	consulConf.Address = dsn
	client, err := consulapi.NewClient(consulConf)
	if err != nil {
		return nil, err
	}

	registrar := &Registrar{
		dsn:         dsn,
		client:      client,
		servicesMap: make(map[string]*grpc.ClientConn),
	}

	return registrar, nil
}

func (r *Registrar) RegisterService(serviceInfo *ServiceInfo) error {
	id := fmt.Sprintf("%s:%d", serviceInfo.Name, serviceInfo.NodeId)

	var srv *consulapi.AgentServiceRegistration
	if serviceInfo.Scheme == "grpc" {
		srv = &consulapi.AgentServiceRegistration{
			ID:      id,               // 服务唯一ID
			Name:    serviceInfo.Name, // 服务名称
			Tags:    []string{serviceInfo.Name},
			Port:    serviceInfo.Port,
			Address: serviceInfo.Address,
			Check: &consulapi.AgentServiceCheck{
				CheckID:                        id,
				GRPC:                           fmt.Sprintf("%s:%d", serviceInfo.Address, serviceInfo.Port),
				Interval:                       "5s",  // 每5秒检测一次
				Timeout:                        "5s",  // 5秒超时
				DeregisterCriticalServiceAfter: "10s", // 超时10秒注销节点
			},
		}
	} else if serviceInfo.Scheme == "http" {
		srv = &consulapi.AgentServiceRegistration{
			ID:      id,               // 服务唯一ID
			Name:    serviceInfo.Name, // 服务名称
			Tags:    []string{serviceInfo.Name},
			Port:    serviceInfo.Port,
			Address: serviceInfo.Address,
			Check: &consulapi.AgentServiceCheck{
				CheckID:                        id,
				HTTP:                           fmt.Sprintf("http://%s:%d/%s", serviceInfo.Address, serviceInfo.Port, "health"),
				Interval:                       "5s",  // 每5秒检测一次
				Timeout:                        "5s",  // 5秒超时
				DeregisterCriticalServiceAfter: "10s", // 超时10秒注销节点
			},
		}
	} else {
		return fmt.Errorf("scheme not support")
	}

	return r.client.Agent().ServiceRegister(srv)
}

func (r *Registrar) UnRegister(serviceInfo *ServiceInfo) error {
	id := fmt.Sprintf("%s:%d", serviceInfo.Name, serviceInfo.NodeId)

	return r.client.Agent().ServiceDeregister(id)
}
