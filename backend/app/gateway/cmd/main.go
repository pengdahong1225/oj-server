package main

import (
	"oj-server/app/gateway/internal/server"
)

func main() {
	s := server.Server{}
	s.Scheme = "http"

	err := s.Init()
	if err != nil {
		panic(err)
	}
	s.Start()
}
