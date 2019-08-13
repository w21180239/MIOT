package main

import "gitlab.com/LICOTEK/DuerOS/server"

func main() {
	server := server.NewServer()

	server.Run(":8080")
}
