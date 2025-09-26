package registry

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

func (r *Registry) getGrpcConnection(name string) (*grpc.ClientConn, error) {
	r.mux.RLock()
	conn, ok := r.servicesMap[name]
	if ok && conn.GetState() == connectivity.Ready {
		return conn, nil
	}
	r.mux.RUnlock()

	return r.createGrpcConnection(name)
}

func (r *Registry) createGrpcConnection(name string) (*grpc.ClientConn, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	// 拿到锁后再次检查
	conn, ok := r.servicesMap[name]
	if ok && conn.GetState() == connectivity.Ready {
		return conn, nil
	}

	// 创建服务连接池
	var err error
	target := fmt.Sprintf("consul://%s/%s?wait=10s&healthy=true", r.addr, name)
	conn, err = grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),                // 不安全连接
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), // 负载均衡，轮训策略
	)
	if err != nil {
		logrus.Errorf("failed to create connection with %s: %v", name, err)
		return nil, err
	}

	r.servicesMap[name] = conn

	return conn, nil
}
