package main

import (
	"fmt"
	"oj-server/global"
	"oj-server/svr/user/internal/configs"
	"oj-server/svr/user/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 加载配置
	configPath := fmt.Sprintf("%s/%s", global.ConfigPath, "user_config.yaml")
	err := configs.LoadServerConfigFile(configPath)
	if err != nil {
		panic(err)
	}
	appCfgPath := fmt.Sprintf("%s/%s", global.ConfigPath, "app_config.yaml")
	if err = configs.LoadAppConfigFile(appCfgPath); err != nil {
		panic(err)
	}

	// 新建服务
	app := server.NewServer()
	if err = app.Init(); err != nil {
		panic(err)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		app.Stop()
		time.Sleep(time.Second)
		os.Exit(0)
	}()

	// 启动服务
	app.Run()
}
