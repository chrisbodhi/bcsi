package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
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
			if cmd == "get" {
				conn, err := net.Dial("tcp", ":8888")
				if err != nil {
					fmt.Println(err)
				}
				defer conn.Close()
				conn.Write([]byte(input))
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(string(buf[:n]))
			} else if cmd == "set" {
				// TODO: validate equals sign
				// TODO: do I care about spaces around the equals sign?
				conn, err := net.Dial("tcp", ":8888")
				if err != nil {
					fmt.Println(err)
				}
				defer conn.Close()
				conn.Write([]byte(input))
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(string(buf[:n]))
			} else {
				fmt.Println("`get` or `set`?")
			}
		case <-sigCh:
			fmt.Println("\nToodles!")
			return
		}
	}
}
