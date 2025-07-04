package main

import (
	"oj-server/app/problem-service/internal"
)

func main() {
	server := internal.Server{}
	server.Name = "problem-service"
	server.SrvType = "grpc"

	err := server.Init()
	if err != nil {
		panic(err)
	}
	server.Start()
}
