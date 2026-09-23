package main

import (
	"oj-server/svr/user/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	app := server.NewServer()
	if err := app.Init(); err != nil {
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
	app.Run()
}
