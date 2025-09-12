package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"oj-server/global"
	"oj-server/module/configManager"
	"oj-server/module/logger"
	"oj-server/module/registry"
	"oj-server/svr/common"
	"oj-server/svr/judge/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var (
		server_config string
		help          bool
		judge_server  common.IServer
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

	judge_server = server.NewServer()
	if err = judge_server.Init(); err != nil {
		logrus.Fatalf("Failed to init server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		logrus.Errorf("Recv signal: %v", sig)
		judge_server.Stop()
		time.Sleep(time.Second)
		os.Exit(0)
	}()

	// 启动
	judge_server.Run()
}
