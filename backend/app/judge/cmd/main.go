package main

import "github.com/pengdahong1225/oj-server/backend/app/judge/internal/server"

func main() {
	s := server.Server{}
	s.Name = "judge-service"
	s.Scheme = "grpc"

	err := s.Init()
	if err != nil {
		panic(err)
	}
	s.Start()
}
