package main

import "github.com/pengdahong1225/oj-server/backend/app/gateway-service/internal"

func main() {
	server := internal.Server{}
	server.Name = "gateway-service"
	server.SrvType = "http"

	err := server.Init()
	if err != nil {
		panic(err)
	}
	server.Start()
}
