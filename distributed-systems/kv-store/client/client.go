package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/chrisbodhi/bcsi/distributed-systems/kv-store/utils"
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
			if !strings.HasPrefix(input, "get") && !strings.HasPrefix(input, "set") {
				fmt.Println("Usage: `get KEY` or `set KEY=VALUE`")
				continue
			}

			args := strings.Split(input, " ")
			cmd := args[0]

			if cmd == "get" {
				sendAndReceive(input)
			} else if cmd == "set" {
				setArg := utils.WithSpace(args[1:])
				// Validate equals sign
				err := utils.ValidateSet(setArg)
				if err != nil {
					fmt.Println(err)
				}
				sendAndReceive(input)
			} else {
				fmt.Println("`get` or `set`?")
			}
		case <-sigCh:
			fmt.Println("\nToodles!")
			return
		}
	}
}

func sendAndReceive(input string) {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	conn.Write([]byte(input))

	buffer := make([]byte, 1024)
	bufferLen, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(buffer[:bufferLen]))
}
