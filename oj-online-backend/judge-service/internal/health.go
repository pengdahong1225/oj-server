package internal

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"judge-service/global"
	"judge-service/utils"
	"net"
)

type HealthServer struct {
}

func (receiver *HealthServer) start() {
	ip, _ := utils.GetOutboundIP()
	// 监听一个tcp端口，用于健康检查
	dsn := fmt.Sprintf("%s:%d", ip, global.ConfigInstance.System_.Port)
	_, err := net.Listen("tcp", dsn)
	if err != nil {
		logrus.Errorf("listen tcp error: %s", err.Error())
		panic(err)
	}

	// 注册
	registry, err := global.NewRegistry()
	if err != nil {
		panic(err)
	}
	if err = registry.RegisterService(global.ConfigInstance.System_.Name, ip.String(), global.ConfigInstance.System_.Port); err != nil {
		panic(err)
	}

	// 阻塞
	select {}
}
