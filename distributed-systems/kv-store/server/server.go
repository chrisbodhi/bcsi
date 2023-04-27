package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/chrisbodhi/bcsi/distributed-systems/kv-store/utils"
)

var mem = make(map[string]string)

var STORAGE_BASE = "storage.json"

func main() {
	table := "default"
	loadDatastore(table)
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Listening on :8888...")
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
	}

	key := string(buf[:n])
	args := strings.Split(key, " ")
	cmd := args[0]

	if cmd == "get" {
		got := Get(args[1])
		conn.Write([]byte(got))
	} else if cmd == "set" {
		err := utils.ValidateSet(args[1])
		if err != nil {
			fmt.Println(err)
			conn.Write([]byte("<validation error>"))
			return
		}
		setPieces := strings.Split(utils.WithSpace(args[1:]), "=")
		k, v := setPieces[0], utils.WithSpace(setPieces[1:])
		Set(k, v)
		conn.Write([]byte("ok"))
	} else {
		fmt.Println("Unknown command:", cmd)
	}
}

func loadDatastore(table string) {
	storage := fmt.Sprintf("%s_%s", table, STORAGE_BASE)
	// Create storage file if it doesn't exist
	if _, err := os.Stat(storage); os.IsNotExist(err) {
		_, err := os.Create(storage)
		if err != nil {
			fmt.Println(err)
		}
	}
	jsonFile, err := os.Open(storage)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(byteValue, &mem)
	fmt.Println(mem)
}

func updateDatastore(table string) {
	storage := fmt.Sprintf("%s_%s", table, STORAGE_BASE)
	// Create storage file if it doesn't exist
	if _, err := os.Stat(storage); os.IsNotExist(err) {
		_, err := os.Create(storage)
		if err != nil {
			fmt.Println(err)
		}
	}
	jsonFile, err := os.OpenFile(storage, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	jsonData, err := json.Marshal(mem)
	if err != nil {
		fmt.Println(err)
	}
	jsonFile.Write(jsonData)
}

func Get(key string) string {
	val, ok := mem[key]
	if !ok {
		return "<not found>"
	}
	return val
}

func Set(key, value string) {
	table := "default"
	mem[key] = value
	// Flush mem to table_storage.json
	updateDatastore(table)
	fmt.Println(key, value)
}
