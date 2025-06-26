package registry

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"google.golang.org/grpc"
	"sync"
)

type Registry struct {
	addr   string
	client *consulapi.Client
	scheme string // http or grpc

	servicesMap map[string]*grpc.ClientConn // 服务连接池
	mux         sync.Mutex
}

func (r *Registry) registerService(info *pb.PBNodeInfo) error {
	id := fmt.Sprintf("%d:%d", info.NodeType, info.NodeId)

	var srv *consulapi.AgentServiceRegistration
	if r.scheme == "grpc" {
		srv = &consulapi.AgentServiceRegistration{
			ID:      id,        // 服务唯一ID
			Name:    info.Name, // 服务名称
			Tags:    []string{info.Name},
			Port:    int(info.Port),
			Address: info.Ip,
			Check: &consulapi.AgentServiceCheck{
				CheckID:                        id,
				GRPC:                           fmt.Sprintf("%s:%d", info.Ip, info.Port),
				Interval:                       "5s",  // 每5秒检测一次
				Timeout:                        "5s",  // 5秒超时
				DeregisterCriticalServiceAfter: "10s", // 超时10秒注销节点
			},
		}
	} else if r.scheme == "http" {
		srv = &consulapi.AgentServiceRegistration{
			ID:      id,        // 服务唯一ID
			Name:    info.Name, // 服务名称
			Tags:    []string{info.Name},
			Port:    int(info.Port),
			Address: info.Ip,
			Check: &consulapi.AgentServiceCheck{
				CheckID:                        id,
				HTTP:                           fmt.Sprintf("http://%s:%d/%s", info.Ip, info.Port, "health"),
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

func (r *Registry) unRegister(info *pb.PBNodeInfo) error {
	id := fmt.Sprintf("%d:%d", info.NodeType, info.NodeId)
	err := r.client.Agent().ServiceDeregister(id)
	if err != nil {
		return err
	}
	return nil
}
