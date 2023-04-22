package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Starting client...")

	inputCh := make(chan string)
	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			inputCh <- scanner.Text()
		}
	}()

	for {
		fmt.Print("🔑 ")
		select {
		case input := <-inputCh:
			fmt.Println("You entered:", input)
		case <-sigCh:
			fmt.Println("\nToodles!")
			return
		}
	}
}
