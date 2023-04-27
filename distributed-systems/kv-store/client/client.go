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

	tables := []string{"default"}
	displayTables := strings.Join(tables, ", ")

	rl, err := readline.New(fmt.Sprintf("🔑 (%s) ", displayTables))
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
		if !strings.HasPrefix(line, "get") && !strings.HasPrefix(line, "pick") && !strings.HasPrefix(line, "set") {
			// TODO: add ability to drop table(s)
			fmt.Println("Usage: `get KEY` or `set KEY=VALUE` or `pick TABLE1 TABLE2 TABLE3`")
			continue
		}

		args := strings.Split(line, " ")
		cmd := args[0]

		if cmd == "pick" {
			tables = args[1:]
			displayTables = strings.Join(tables, ", ")
			fmt.Println("Using tables:", displayTables)
			// Reload readline with new prompt
			rl, err = readline.New(fmt.Sprintf("🔑 (%s) ", displayTables))
			if err != nil {
				fmt.Println(err)
				break
			}
		} else if cmd == "get" {
			sendAndReceive(line, tables)
		} else if cmd == "set" {
			setArg := utils.WithSpace(args[1:])
			// Validate equals sign
			err := utils.ValidateSet(setArg)
			if err != nil {
				fmt.Println(err)
			}
			sendAndReceive(line, tables)
		} else {
			fmt.Println("`get` or `set`?")
		}
	}
}

func sendAndReceive(input string, tables []string) {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	for _, table := range tables {
		// TODO: verify table exists; we want to handle typos in the table name
		//       so that we don't end up with "default" and "deafult" and "defualt" &c.
		fmt.Println("Using table:", table)
		send := fmt.Sprintf("%s %s", table, input)
		conn.Write([]byte(send))

		buffer := make([]byte, 1024)
		bufferLen, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(buffer[:bufferLen]))
	}

}
