package main

import (
	"fmt"

	"github.com/chrisbodhi/bcsi/ds-for-storage-and-retrieval/01/goldb"
)

func main() {
	fmt.Println("Getting started...")
	db := goldb.InitDb()
	defer goldb.CloseDb(db)
	fmt.Println("Opened", db)
	err := goldb.Put(db, "foo", "bar")
	if err != nil {
		fmt.Println(err)
	}
	s := goldb.Get(db, "foo")
	fmt.Println("Hello,", s)
}
