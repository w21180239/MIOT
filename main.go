package main

import "server"

func main() {
	server := server.NewServer()

	server.Run(":8080")
}
