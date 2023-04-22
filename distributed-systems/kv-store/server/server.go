package server

import "fmt"

var mem = make(map[string]string)

func Get(key string) {
	// TODO: error handling
	fmt.Println(mem[key])
}

func Set(key, value string) {
	// TODO: error handling
	mem[key] = value
	fmt.Println("set", key, value, "successfully")
}
