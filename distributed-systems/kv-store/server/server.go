package server

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var mem = make(map[string]string)

var STORAGE = "storage.json"

func init() {
	loadDatastore()
	fmt.Println(fmt.Sprintf("Successfully opened %s\nand loaded it into memory:", STORAGE), mem)
}

func loadDatastore() {
	jsonFile, err := os.Open(STORAGE)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
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

func Get(key string) {
	loadDatastore()
	// TODO: error handling
	fmt.Println(mem[key])
}

func Set(key, value string) {
	// TODO: error handling
	mem[key] = value
	// Flush mem to storage.json
	updateDatastore()
	fmt.Println(key, value)
}
