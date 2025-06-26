package main

import (
	"github.com/pengdahong1225/oj-server/backend/app/gateway/internal/server"
)

func main() {
	s := server.Server{}
	s.Name = "gateway-service"
	s.Scheme = "http"

	err := s.Init()
	if err != nil {
		panic(err)
	}
	s.Start()
}
