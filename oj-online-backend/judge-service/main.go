package main

import (
	"judge-service/server"
)

func main() {
	server := server.Server{}
	server.Start()
}
