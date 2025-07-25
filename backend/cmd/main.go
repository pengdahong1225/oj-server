package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"oj-server/app"
	"oj-server/global"
	"oj-server/module/configManager"
	"oj-server/module/logger"
	"oj-server/module/registry"
	"os"
)

func main() {
	var (
		server_config string
		help          bool
	)
	flag.StringVar(&server_config, "f", "server_config.yaml", "server config")
	flag.BoolVar(&help, "h", false, "help")
	flag.Parse()

	actualArgs := len(os.Args[1:])
	if actualArgs < 1 || help {
		flag.Usage()
		os.Exit(1)
	}

	// 加载配置
	configPath := fmt.Sprintf("%s/%s", global.ConfigPath, server_config)
	err := configManager.LoadServerConfigFile(configPath)
	if err != nil {
		panic(err)
	}
	appCfgPath := fmt.Sprintf("%s/%s", global.ConfigPath, "app_config.yaml")
	err = configManager.LoadAppConfigFile(appCfgPath)
	if err != nil {
		panic(err)
	}

	// 初始化日志
	serverCfg := configManager.ServerConf
	err = logger.Init(global.LogPath, serverCfg.NodeType, logrus.DebugLevel)
	if err != nil {
		panic(err)
	}

	// 初始化注册中心
	err = registry.Init()
	if err != nil {
		panic(err)
	}
	logrus.Debugf("--------------- node_type:%v, node_id:%v, host:%v, port:%v, scheme:%v ---------------", serverCfg.NodeType, serverCfg.NodeId, serverCfg.Host, serverCfg.Port, serverCfg.Scheme)

	server, err := app.NewServer()
	if err != nil {
		logrus.Fatalf("Failed to create server: %v", err)
	}
	err = server.Init()
	if err != nil {
		logrus.Fatalf("Failed to init server: %v", err)
	}
	server.Run()
}
