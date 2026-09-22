package registry

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	consulapi "github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConsulRegistry struct {
	dsn    string
	client *consulapi.Client

	connMu sync.Mutex
	conns  map[string]*grpc.ClientConn

	sf singleflight.Group // 防止并发 Dial 同一服务重复建连

	regMu      sync.Mutex
	registered map[string]struct{}
}

var _ Registry = (*ConsulRegistry)(nil)

func NewConsulRegistry(dsn string) (*ConsulRegistry, error) {
	consulConf := consulapi.DefaultConfig()
	consulConf.Address = dsn
	client, err := consulapi.NewClient(consulConf)
	if err != nil {
		return nil, err
	}
	return &ConsulRegistry{
		dsn:        dsn,
		client:     client,
		conns:      make(map[string]*grpc.ClientConn),
		registered: make(map[string]struct{}),
	}, nil
}

func (r *ConsulRegistry) Register(ctx context.Context, inst ServiceInfo) (func() error, error) {
	id := inst.ServiceID()

	r.regMu.Lock()
	if _, ok := r.registered[id]; ok {
		r.regMu.Unlock()
		return nil, fmt.Errorf("service %s already registered", id)
	}
	r.regMu.Unlock()

	var check *consulapi.AgentServiceCheck
	switch inst.Scheme {
	case "grpc":
		check = &consulapi.AgentServiceCheck{
			CheckID:                        id,
			GRPC:                           fmt.Sprintf("%s:%d", inst.Address, inst.Port),
			Interval:                       "5s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "10s",
		}
	case "http":
		check = &consulapi.AgentServiceCheck{
			CheckID:                        id,
			HTTP:                           fmt.Sprintf("http://%s:%d/health", inst.Address, inst.Port),
			Interval:                       "5s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "10s",
		}
	default:
		return nil, fmt.Errorf("scheme not support: %s", inst.Scheme)
	}

	srv := &consulapi.AgentServiceRegistration{
		ID:      id,
		Name:    inst.Name,
		Tags:    []string{inst.Name, inst.Scheme},
		Port:    inst.Port,
		Address: inst.Address,
		Check:   check,
	}

	if err := r.client.Agent().ServiceRegister(srv); err != nil {
		return nil, fmt.Errorf("consul register: %w", err)
	}

	r.regMu.Lock()
	r.registered[id] = struct{}{}
	r.regMu.Unlock()

	return func() error {
		return r.Deregister(context.Background(), id)
	}, nil
}

func (r *ConsulRegistry) Deregister(ctx context.Context, serviceID string) error {
	r.regMu.Lock()
	if _, ok := r.registered[serviceID]; !ok {
		r.regMu.Unlock()
		return nil // 幂等
	}
	r.regMu.Unlock()

	if err := r.client.Agent().ServiceDeregister(serviceID); err != nil {
		return fmt.Errorf("consul deregister: %w", err)
	}

	r.regMu.Lock()
	delete(r.registered, serviceID)
	r.regMu.Unlock()
	return nil
}

func (r *ConsulRegistry) Discover(ctx context.Context, serviceName string) ([]ServiceInfo, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("consul health service: %w", err)
	}

	instances := make([]ServiceInfo, 0, len(entries))
	for _, e := range entries {
		addr := e.Service.Address
		if addr == "" {
			addr = e.Node.Address
		}
		instances = append(instances, ServiceInfo{
			Name:    e.Service.Service,
			Address: addr,
			Port:    e.Service.Port,
			Scheme:  schemeFromTags(e.Service.Tags),
		})
	}
	return instances, nil
}

func (r *ConsulRegistry) Resolve(ctx context.Context, serviceName string) (ServiceInfo, error) {
	instances, err := r.Discover(ctx, serviceName)
	if err != nil {
		return ServiceInfo{}, err
	}
	if len(instances) == 0 {
		return ServiceInfo{}, fmt.Errorf("no healthy instance for service %s", serviceName)
	}
	return instances[rand.Intn(len(instances))], nil
}

func (r *ConsulRegistry) Dial(ctx context.Context, name string) (*grpc.ClientConn, error) {
	// 快路径：已缓存
	r.connMu.Lock()
	if conn, ok := r.conns[name]; ok {
		r.connMu.Unlock()
		return conn, nil
	}
	r.connMu.Unlock()

	// 慢路径：singleflight 保证同一服务只建一次
	v, err, _ := r.sf.Do(name, func() (any, error) {
		r.connMu.Lock()
		if conn, ok := r.conns[name]; ok {
			r.connMu.Unlock()
			return conn, nil
		}
		r.connMu.Unlock()

		conn, err := grpc.NewClient(
			fmt.Sprintf("consul://%s/%s?wait=10s&healthy=true", r.dsn, name),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		)
		if err != nil {
			return nil, err
		}

		r.connMu.Lock()
		r.conns[name] = conn
		r.connMu.Unlock()
		return conn, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*grpc.ClientConn), nil
}

func (r *ConsulRegistry) Close() error {
	r.connMu.Lock()
	defer r.connMu.Unlock()
	for _, conn := range r.conns {
		conn.Close()
	}
	r.conns = make(map[string]*grpc.ClientConn)
	return nil
}

func schemeFromTags(tags []string) string {
	for _, t := range tags {
		if t == "grpc" || t == "http" {
			return t
		}
	}
	return ""
}
