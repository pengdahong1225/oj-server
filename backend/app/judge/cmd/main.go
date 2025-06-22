package main

import (
	"github.com/pengdahong1225/oj-server/backend/app/judge-service/internal"
)

func main() {
	server := internal.Server{}
	server.Name = "judge-service"
	server.SrvType = "http"

	err := server.Init()
	if err != nil {
		panic(err)
	}
	server.Start()
}
