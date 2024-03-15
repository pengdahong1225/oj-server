package main

import "judge-service/internal"

func main() {
	server := internal.Server{}
	server.Start()
}
