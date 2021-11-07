package goldb

// #cgo LDFLAGS: -lleveldb
// #include "leveldb/c.h"
// #include "level.h"
import "C"
import (
	"fmt"
	"unsafe"
)

type DB struct {
	db *C.leveldb_t
}

func InitDb() *C.leveldb_t {
	var db *C.leveldb_t
	db = C.initlevel()
	return db
}

func CloseDb(db *C.leveldb_t) {
	C.closelevel(db)
}

func Put(db *C.leveldb_t, key string, value string) error {
	keyBytes := []byte(key)
	valueBytes := []byte(value)
	keyStart := (*C.char)(unsafe.Pointer(&keyBytes[0]))
	valueStart := (*C.char)(unsafe.Pointer(&valueBytes[0]))
	returnVal := C.putlevel(db, keyStart, valueStart)
	if returnVal == 0 {
		return nil
	} else {
		return fmt.Errorf("Error")
	}
}

func Get(db *C.leveldb_t, key string) string {
	var value *C.char
	keyBytes := []byte(key)
	fmt.Println(keyBytes)
	keyStart := (*C.char)(unsafe.Pointer(&keyBytes[0]))
	value = C.getlevel(db, keyStart)
	return C.GoString(value)
}

func Empty() {
	fmt.Println("Empty")
}
