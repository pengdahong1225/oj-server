package main

import (
	"oj-server/app/user/internal/server"
)

func main() {
	s := server.Server{}
	s.Name = "user-service"
	s.Scheme = "grpc"

	err := s.Init()
	if err != nil {
		panic(err)
	}
	s.Start()
}
