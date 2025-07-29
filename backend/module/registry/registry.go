package registry

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"oj-server/module/configManager"
	"sync"
)

type Registry struct {
	addr   string
	client *consulapi.Client
	scheme string // http or grpc

	servicesMap map[string]*grpc.ClientConn // 服务连接池
	mux         sync.Mutex
}

func (r *Registry) registerService() error {
	cfg := configManager.ServerConf

	logrus.Debugf("注册服务: %s:%d", cfg.Host, cfg.Port)

	id := fmt.Sprintf("%s:%d", cfg.NodeType, cfg.NodeId)

	var srv *consulapi.AgentServiceRegistration
	if r.scheme == "grpc" {
		srv = &consulapi.AgentServiceRegistration{
			ID:      id,           // 服务唯一ID
			Name:    cfg.NodeType, // 服务名称
			Tags:    []string{cfg.NodeType},
			Port:    cfg.Port,
			Address: cfg.Host,
			Check: &consulapi.AgentServiceCheck{
				CheckID:                        id,
				GRPC:                           fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
				Interval:                       "5s",  // 每5秒检测一次
				Timeout:                        "5s",  // 5秒超时
				DeregisterCriticalServiceAfter: "10s", // 超时10秒注销节点
			},
		}
	} else if r.scheme == "http" {
		srv = &consulapi.AgentServiceRegistration{
			ID:      id,           // 服务唯一ID
			Name:    cfg.NodeType, // 服务名称
			Tags:    []string{cfg.NodeType},
			Port:    cfg.Port,
			Address: cfg.Host,
			Check: &consulapi.AgentServiceCheck{
				CheckID:                        id,
				HTTP:                           fmt.Sprintf("http://%s:%d/%s", cfg.Host, cfg.Port, "health"),
				Interval:                       "5s",  // 每5秒检测一次
				Timeout:                        "5s",  // 5秒超时
				DeregisterCriticalServiceAfter: "10s", // 超时10秒注销节点
			},
		}
	} else {
		return fmt.Errorf("scheme not support")
	}

	err := r.client.Agent().ServiceRegister(srv)
	if err != nil {
		return err
	}
	return nil
}

func (r *Registry) unRegister() error {
	cfg := configManager.ServerConf
	id := fmt.Sprintf("%s:%d", cfg.NodeType, cfg.NodeId)
	err := r.client.Agent().ServiceDeregister(id)
	if err != nil {
		return err
	}
	return nil
}
