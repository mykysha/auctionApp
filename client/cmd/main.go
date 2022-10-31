package main

import (
	"sync"

	"github.com/nndergunov/auctuionApp/client/pkg/client"
)

func main() {
	clientService := client.NewClient("ws://localhost:7000/ws", "http://localhost:7000/lp")

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		defer wg.Done()
		clientService.WSConnect()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		clientService.LongPollConnect()
	}()
}
