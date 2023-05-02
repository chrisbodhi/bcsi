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
		if !strings.HasPrefix(line, "get") && !strings.HasPrefix(line, "pick") && !strings.HasPrefix(line, "set") && !strings.HasPrefix(line, "drop") {
			fmt.Println("Usage: `get KEY` or `set KEY=VALUE` or `pick TABLE1 TABLE2 TABLE3` or `drop TABLE1 TABLE2`")
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
		} else if cmd == "drop" {
			// TODO: do not drop table if it is current table
			//       OR switch to default and then drop
			table := args[1:][0]
			fmt.Println("table to delete", table)
			sendAndReceive(cmd, table)
		} else if cmd == "get" {
			for _, table := range tables {
				sendAndReceive(line, table)
			}
		} else if cmd == "set" {
			err := utils.ValidateSet(args[1:])
			if err != nil {
				fmt.Println(err)
			}
			for _, table := range tables {
				sendAndReceive(line, table)
			}
		} else {
			fmt.Println("`get` or `set`? Maybe `pick` or `drop`?")
		}
	}
}

func sendAndReceive(input string, table string) {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

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
