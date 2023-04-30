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

var mem = make(map[string]map[string][]byte)

var STORAGE_BASE = "storage.json"

func main() {
	// TODO: are these two lines necessary?
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
	displayTables := args[0]
	tables := strings.Split(displayTables, ",")
	cmd := args[1]

	if cmd == "drop" {
		dropped := Drop(displayTables)
		conn.Write([]byte(dropped))
	} else if cmd == "get" {
		got := Get(args[2], tables)
		conn.Write([]byte(fmt.Sprintf("%v", got)))
	} else if cmd == "set" {
		err := utils.ValidateSet(args[2:])
		if err != nil {
			fmt.Println(err)
			conn.Write([]byte("<validation error>"))
			return
		}
		setPieces := utils.InputToSetPieces(args[3:])
		k, v := args[2], setPieces
		Set(k, v, tables)
		conn.Write([]byte(fmt.Sprintf("%v", v)))
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
	localTable := make(map[string][]byte)
	json.Unmarshal(byteValue, &localTable)
	mem[table] = localTable
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
	jsonData, err := json.Marshal(mem[table])
	if err != nil {
		fmt.Println(err)
	}
	jsonFile.Write(jsonData)
}

func Get(key string, tables []string) utils.UserRecord {
	// TODO: placeholder code is just a placeholder
	// This needs to be more intelligent than just
	// printing all values, table by table.
	var last utils.UserRecord
	for _, table := range tables {
		loadDatastore(table)
		val, ok := mem[table][key]
		last = utils.Decode(val)
		if !ok {
			fmt.Printf("Key %s not found in table %s\n", key, table)
		}
	}
	return last
}

func Set(key string, value utils.UserRecord, tables []string) {
	// Flush mem to {table}_storage.json
	for _, table := range tables {
		if _, ok := mem[table]; !ok {
			mem[table] = make(map[string][]byte)
		}
		mem[table][key] = utils.Encode(value)
		updateDatastore(table)
	}
	fmt.Printf("Set %s to %s", key, value)
}

func Drop(table string) string {
	// Remove from mem
	if _, ok := mem[table]; ok {
		delete(mem, table)
	} else {
		return fmt.Sprintf("%s does not exist", table)
	}
	// Rename backing datastore/file
	storage := fmt.Sprintf("%s_%s", table, STORAGE_BASE)
	if err := os.Rename(storage, fmt.Sprintf("dropped_%s", table)); err != nil {
		return fmt.Sprintf("Failed to removing backing datastore for %s", table)
	}

	return fmt.Sprintf("Removed %s", table)
}
