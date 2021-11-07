package goldb

// #cgo LDFLAGS: -lleveldb
// #include "leveldb/c.h"
// #include "level.h"
import "C"
import "fmt"

func Goldb() *C.leveldb_t {
	var db *C.leveldb_t
	db = C.initlevel()
	fmt.Println("Hola, mundo!")
	return db
}

func Empty() {
	fmt.Println("Empty")
}
