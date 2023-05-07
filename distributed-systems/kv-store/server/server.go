package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/chrisbodhi/bcsi/distributed-systems/kv-store/utils"
)

var mem = make(map[string]map[string][]byte)

var STORAGE_BASE = "storage.json"

func Start(port string) {
	// TODO: are these two lines necessary?
	table := "default"
	loadDatastore(table)
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
		dropped, err := Drop(displayTables)
		if err != nil {
			fmt.Println(err)
			msg := fmt.Sprintf("<table error: %s not found>", displayTables)
			conn.Write([]byte(msg))
			return
		}
		utils.WriteAhead(key)
		conn.Write([]byte(dropped))
	} else if cmd == "get" {
		if len(tables) > 1 {
			fmt.Println("Only one table allowed for GET")
			conn.Write([]byte("<validation error: only one table allowed at a time>"))
			return
		}
		got, err := Get(args[2], tables[0])
		if err != nil {
			fmt.Println(err)
			msg := fmt.Sprintf("<%v for %s in %s>", err, args[2], tables[0])
			conn.Write([]byte(msg))
			return
		}
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
		utils.WriteAhead(key)
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

func Get(key string, table string) (utils.UserRecord, error) {
	var last utils.UserRecord
	loadDatastore(table)
	val, ok := mem[table][key]
	if !ok {
		fmt.Printf("Key %s not found in table %s\n", key, table)
		return last, errors.New("Key not found")
	}
	last = utils.Decode(val)
	return last, nil
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

func Drop(table string) (string, error) {
	// Remove from mem
	if _, ok := mem[table]; ok {
		delete(mem, table)
	} else {
		msg := fmt.Sprintf("%s does not exist", table)
		return "", errors.New(msg)
	}
	// Rename backing datastore/file
	storage := fmt.Sprintf("%s_%s", table, STORAGE_BASE)
	if err := os.Rename(storage, fmt.Sprintf("dropped_%s", table)); err != nil {
		msg := fmt.Sprintf("Failed to removing backing datastore for %s", table)
		return msg, nil
	}

	msg := fmt.Sprintf("Removed %s", table)
	return msg, nil
}

func Close() {
	fmt.Println("Closing...")
}
