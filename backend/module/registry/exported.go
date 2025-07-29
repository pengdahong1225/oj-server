package registry

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"oj-server/module/configManager"
)

var (
	instance *Registry
)

func Init() error {
	instance = new(Registry)

	// 配置中心地址
	server_cfg := configManager.ServerConf
	app_cfg := configManager.AppConf.RegistryCfg
	dsn := fmt.Sprintf("%s:%d", app_cfg.Host, app_cfg.Port)
	consulConf := consulapi.DefaultConfig()
	consulConf.Address = dsn
	// client
	c, err := consulapi.NewClient(consulConf)
	if err != nil {
		return err
	}
	instance.addr = dsn
	instance.client = c
	instance.scheme = server_cfg.Scheme

	instance.servicesMap = make(map[string]*grpc.ClientConn)

	return nil
}

func RegisterService() error {
	return instance.registerService()
}
func DeregisterService() error {
	return instance.unRegister()
}
func GetGrpcConnection(name string) (*grpc.ClientConn, error) {
	return instance.getGrpcConnection(name)
}
