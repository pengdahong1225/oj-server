package main

import "db-service/internal"

func main() {
	server := internal.Server{}
	server.Start()
}
