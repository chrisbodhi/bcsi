package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/chrisbodhi/bcsi/distributed-systems/kv-store/utils"
	"github.com/chzyer/readline"
)

func main() {
	fmt.Println("Starting client...")

	rl, err := readline.New("🔑 ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			fmt.Println(err)
			break
		}
		// TODO: improve variable names
		if !strings.HasPrefix(line, "get") && !strings.HasPrefix(line, "set") {
			fmt.Println("Usage: `get KEY` or `set KEY=VALUE`")
			continue
		}

		args := strings.Split(line, " ")
		cmd := args[0]

		if cmd == "get" {
			sendAndReceive(line)
		} else if cmd == "set" {
			setArg := utils.WithSpace(args[1:])
			// Validate equals sign
			err := utils.ValidateSet(setArg)
			if err != nil {
				fmt.Println(err)
			}
			sendAndReceive(line)
		} else {
			fmt.Println("`get` or `set`?")
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
