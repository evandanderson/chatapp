package main

import (
	"chatapp/server/consts"
	"chatapp/server/server"
	"log"
)

func main() {
	server := server.NewServer()
	err := server.Listen(consts.Network, consts.DefaultPort)
	if err != nil {
		log.Fatal(err)
	}
}
