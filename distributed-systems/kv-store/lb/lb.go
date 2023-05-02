package main

import (
	"fmt"
	"sync"

	"github.com/chrisbodhi/bcsi/distributed-systems/kv-store/server"
)

var PORTS = []string{":8889", ":8890"}

func main() {
	fmt.Println("Starting servers...")

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		server.Start(PORTS[0])
		defer wg.Done()
	}()

	go func() {
		server.Start(PORTS[1])
		defer wg.Done()
	}()

	wg.Wait()
}
