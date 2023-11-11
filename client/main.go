package main

import (
	"chatapp/client/client"
	"chatapp/client/consts"
	"log"
	"sync"
)

func main() {
	var wg = &sync.WaitGroup{}
	client := client.NewClient()
	err := client.Dial(consts.Network, ":8080")
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(2)
	go client.HandleInput()
	go client.Receive()
	wg.Wait()
}
