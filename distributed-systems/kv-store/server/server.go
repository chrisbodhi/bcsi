package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var mem = make(map[string]string)

var STORAGE = "storage.json"

func main() {
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
		k, v := strings.Split(args[1], "=")[0], strings.Split(args[1], "=")[1]
		Set(k, v)
		conn.Write([]byte("ok"))
	} else {
		fmt.Println("Unknown command:", cmd)
	}
}

func loadDatastore() {
	jsonFile, err := os.Open(STORAGE)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(byteValue, &mem)
}

func updateDatastore() {
	jsonFile, err := os.OpenFile(STORAGE, os.O_RDWR, 0644)
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
	loadDatastore()
	val, ok := mem[key]
	if !ok {
		return "<not found>"
	}
	return val
}

func Set(key, value string) {
	// TODO: error handling
	mem[key] = value
	// Flush mem to storage.json
	updateDatastore()
	fmt.Println(key, value)
}
