package main

import "github.com/pengdahong1225/oj-server/backend/app/user/internal"

func main() {
	server := internal.Server{}
	server.Name = "user-service"
	server.SrvType = "grpc"

	err := server.Init()
	if err != nil {
		panic(err)
	}
	server.Start()
}
