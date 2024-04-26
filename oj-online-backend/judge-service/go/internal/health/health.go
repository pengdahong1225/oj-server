package health

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"judge-service/services/registry"
	"judge-service/settings"
	"judge-service/utils"
	"net"
)

type HealthServer struct {
}

func (receiver *HealthServer) Loop() {
	ip, _ := utils.GetOutboundIP()
	// 监听一个tcp端口，用于健康检查
	dsn := fmt.Sprintf("%s:%d", ip, settings.Conf.SystemConfig.Port)
	_, err := net.Listen("tcp", dsn)
	if err != nil {
		logrus.Errorf("listen tcp error: %s", err.Error())
		panic(err)
	}

	// 注册
	register, err := registry.NewRegistry(settings.Conf.RegistryConfig)
	if err != nil {
		panic(err)
	}
	if err = register.RegisterService(settings.Conf.SystemConfig.Name, ip.String(), settings.Conf.SystemConfig.Port); err != nil {
		panic(err)
	}

	// 阻塞
	select {}
}
