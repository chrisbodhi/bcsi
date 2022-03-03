package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

type trivialBloomFilter struct {
	data []uint64
}

func newTrivialBloomFilter() *trivialBloomFilter {
	return &trivialBloomFilter{
		data: make([]uint64, 1000),
	}
}

func (b *trivialBloomFilter) add(item string) {}

func (b *trivialBloomFilter) maybeContains(item string) bool {
	// Technically, any item "might" be in the set
	return true
}

func (b *trivialBloomFilter) memoryUsage() int {
	return binary.Size(b.data)
}

type bloomFilter interface {
	add(item string)

	// `false` means the item is definitely not in the set
	// `true` means the item might be in the set
	maybeContains(item string) bool

	// Number of bytes used in any underlying storage
	memoryUsage() int
}

type cbBloomFilter struct {
	data    big.Int
	hashFns int
}

// newCbBloomFilter takes arguments that determine
// characteristics of the underlying data structure,
// e.g. the bitset, in order to allow us to test various
// properties to find the best parameters for our
// circumstances.
// º p is the probability of false positives that
// we are willing to accept.
// º cap is the capacity of our Bloom filter, so that
// we can contain any set `S` containing up to `cap`
// elements.
func newCbBloomFilter(p float64, cap int) *cbBloomFilter {
	k := math.Log(1 / p)
	// k := 6.0
	fmt.Println("p:", p)
	fmt.Println("cap:", cap)
	fmt.Println("k:", k)
	// size-1 also becomes the max returned value of
	// our hash functions
	// TODO revisit this calculation; should I still have
	// to multiply by 10 to get a good result?
	sizeF := (float64(cap) / math.Log(2)) * math.Log2(1/p)
	size := sizeF * 10.0
	fmt.Println("size:", size)
	return &cbBloomFilter{
		data:    *big.NewInt(int64(size)), // TODO is this zero'd?
		hashFns: int(k),
	}
}

func (cb *cbBloomFilter) add(item string) {
	for i := 0; i < cb.hashFns; i++ {
		index := Hashing(item, i, int(cb.data.Int64()))
		cb.data.SetBit(&cb.data, int(index), 1)
	}
}

func (cb *cbBloomFilter) maybeContains(item string) bool {
	checks := make([]bool, cb.hashFns)
	for i := 0; i < cb.hashFns; i++ {
		index := Hashing(item, i, int(cb.data.Int64()))
		checks[i] = cb.Check(int(index), &cb.data)
	}
	for _, check := range checks {
		if !check {
			return check
		}
	}
	return true // all checks came back `true`
}

func (cb *cbBloomFilter) memoryUsage() int {
	return binary.Size(cb.data)
}
