package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	server "github.com/chrisbodhi/bcsi/distributed-systems/kv-store/server"
)

func main() {
	fmt.Println("Starting client...")

	inputCh := make(chan string)
	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// TODO: replace with readline -- https://pkg.go.dev/github.com/chzyer/readline
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			inputCh <- scanner.Text()
		}
	}()

	for {
		fmt.Print("🔑 ")
		select {
		case input := <-inputCh:
			// TODO: improve variable names
			args := strings.Split(input, " ")
			if len(args) != 2 {
				fmt.Println("Usage: `get KEY` or `set KEY=VALUE`")
				continue
			}
			cmd := args[0]
			rest := args[1]
			if cmd == "get" {
				server.Get(rest)
			} else if cmd == "set" {
				// TODO: validate equals sign
				// TODO: do I care about spaces around the equals sign?
				k, v := strings.Split(rest, "=")[0], strings.Split(rest, "=")[1]
				// TODO: do I want to consider handling types?
				server.Set(k, v)
			} else {
				fmt.Println("`get` or `set`?")
			}
		case <-sigCh:
			fmt.Println("\nToodles!")
			return
		}
	}
}