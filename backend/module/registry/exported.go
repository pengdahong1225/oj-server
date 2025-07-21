package registry

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"oj-server/module/settings"
	"oj-server/proto/pb"
)

var (
	instance *Registry
)

func Init(scheme string) error {
	instance = new(Registry)

	// 配置中心地址
	cfg := settings.AppConf.RegistryCfg
	dsn := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	consulConf := consulapi.DefaultConfig()
	consulConf.Address = dsn
	// client
	c, err := consulapi.NewClient(consulConf)
	if err != nil {
		return err
	}
	instance.addr = dsn
	instance.client = c
	instance.scheme = scheme

	return nil
}

func RegisterService(info *pb.PBNodeInfo) error {
	return instance.registerService(info)
}
func DeregisterService(info *pb.PBNodeInfo) error {
	return instance.unRegister(info)
}
func GetGrpcConnection(name string) (*grpc.ClientConn, error) {
	return instance.getGrpcConnection(name)
}
