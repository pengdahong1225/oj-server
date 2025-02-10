package main

import "github.com/pengdahong1225/oj-server/backend/app/daemon-service/internal"

func main() {
	server := internal.Server{}
	server.Name = "daemon-service"
	server.SrvType = "http"

	err := server.Init()
	if err != nil {
		panic(err)
	}
	server.Start()
}
