package main

// #cgo LDFLAGS: -lleveldb
// #include "leveldb/c.h"
import "C"

import (
	"fmt"
)

func goldb() *C.leveldb_t {
	var db *C.leveldb_t
	var dbOptions *C.leveldb_options_t
	name := C.CString("temp")
	var err *C.char
	db = C.leveldb_open(dbOptions, name, &err)
	// defer C.free(unsafe.Pointer(name))
	// defer C.free(unsafe.Pointer(err))
	fmt.Println("Hola, mundo!")
	return db
}

func main() {
	goldb()
}
