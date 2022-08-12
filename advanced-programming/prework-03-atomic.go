package main

import (
	"fmt"
	"sync/atomic"
)

type idService interface {
	// Returns values in ascending order; it should be safe to call
	// getNext() concurrently without any additional synchronization.
	getNext() uint64
}

type atomiclike struct {
	current uint64
}

func (a *atomiclike) getNext() uint64 {
	pac := &a.current
	atomic.StoreUint64(pac, atomic.AddUint64(pac, uint64(1)))
	return *pac
}

func main() {
	fmt.Println("atomically!")
	like := atomiclike{uint64(9)}
	for i := 0; i < 10; i++ {
		fmt.Println(like.getNext())
	}
}
