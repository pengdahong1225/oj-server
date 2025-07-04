package registry

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func (r *Registry) getGrpcConnection(name string) (*grpc.ClientConn, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	conn, ok := r.servicesMap[name]
	if ok {
		return conn, nil
	}
	// 创建服务连接池
	// 加载证书
	var (
		creds credentials.TransportCredentials
		err   error
	)
	creds, err = credentials.NewClientTLSFromFile("cert/server.pem", "")
	if err != nil {
		logrus.Errorf("Failed to create TLS credentials %v", err)
		return nil, err
	}
	target := fmt.Sprintf("consul://%s/%s?wait=10s&healthy=true", r.addr, name)
	conn, err = grpc.NewClient(
		target,
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), // 负载均衡，轮训策略
	)
	if err != nil {
		logrus.Errorf("Failed to create connection with %s: %v", name, err)
		return nil, err
	}
	// 记录
	r.servicesMap[name] = conn
	return conn, nil
}
