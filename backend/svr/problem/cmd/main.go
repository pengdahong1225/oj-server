package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"oj-server/global"
	"oj-server/module/configs"
	"oj-server/module/logger"
	"oj-server/module/registry"
	"oj-server/svr/problem/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 加载配置
	configPath := fmt.Sprintf("%s/%s", global.ConfigPath, "server_config.yaml")
	err := configs.LoadServerConfigFile(configPath)
	if err != nil {
		panic(err)
	}
	appCfgPath := fmt.Sprintf("%s/%s", global.ConfigPath, "app_config.yaml")
	err = configs.LoadAppConfigFile(appCfgPath)
	if err != nil {
		panic(err)
	}

	// 初始化日志
	serverCfg := configs.ServerConf
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

	app := server.NewServer()
	if err = app.Init(); err != nil {
		logrus.Fatalf("Failed to init server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		logrus.Errorf("Recv signal: %v", sig)
		app.Stop()
		time.Sleep(time.Second)
		os.Exit(0)
	}()

	// 启动
	app.Run()
}
