package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/chrisbodhi/bcsi/distributed-systems/kv-store/utils"
)

func main() {
	// port variable comes from first command line argument
	port := fmt.Sprintf(":%s", os.Args[1])
	// TODO: are these two lines necessary?
	table := "default"
	utils.LoadDatastore(port, table)
	fmt.Println("received port", port)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("!!", err)
	}
	fmt.Printf("Listening on %s...\n", port)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("err with", port, ":", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	addrComponents := strings.Split(conn.LocalAddr().String(), ":")
	portNumber := addrComponents[len(addrComponents)-1]
	port := fmt.Sprintf(":%s", portNumber)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
	}

	key := string(buf[:n])
	args := strings.Split(key, " ")
	displayTables := args[0]
	tables := strings.Split(displayTables, ",")
	cmd := args[1]

	switch cmd {
	case "drop":
		dropped, err := Drop(displayTables, port)
		if err != nil {
			fmt.Println(err)
			msg := fmt.Sprintf("<table error: %s not found>", displayTables)
			conn.Write([]byte(msg))
			return
		}
		// TODO: make more robust, since conn.Write may finish first
		go utils.WriteAhead(key)
		conn.Write([]byte(dropped))
	case "get":
		if len(tables) > 1 {
			fmt.Println("Only one table allowed for GET")
			conn.Write([]byte("<validation error: only one table allowed at a time>"))
			return
		}
		got, err := Get(args[2], tables[0], port)
		if err != nil {
			fmt.Println(err)
			msg := fmt.Sprintf("<%v for %s in %s>", err, args[2], tables[0])
			conn.Write([]byte(msg))
			return
		}
		conn.Write([]byte(fmt.Sprintf("%v", got)))
	case "set":
		err := utils.ValidateSet(args[2:])
		if err != nil {
			fmt.Println(err)
			conn.Write([]byte("<validation error>"))
			return
		}
		setPieces := utils.InputToSetPieces(args[3:])
		k, v := args[2], setPieces
		// TODO: make more robust, since Set and conn.Write may finish first
		go utils.WriteAhead(key)
		Set(k, v, tables, port)
		conn.Write([]byte(fmt.Sprintf("%v", v)))
	default:
		fmt.Println("Unknown command:", cmd)
	}
}

func Get(key, table, port string) (utils.UserRecord, error) {
	var last utils.UserRecord
	utils.LoadDatastore(port, table)
	val, ok := utils.Get(key, table)
	if !ok {
		fmt.Printf("Key %s not found in table %s\n", key, table)
		return last, errors.New("key not found")
	}
	last = utils.Decode(val)
	return last, nil
}

func Set(key string, value utils.UserRecord, tables []string, port string) {
	for _, table := range tables {
		utils.Set(key, value, table, port)
		fmt.Printf("Set %s to %s in %s\n", key, value, table)
	}
}

func Drop(table, port string) (string, error) {
	msg, err := utils.Drop(table, port)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return msg, nil
}
