package main

import (
	"hash/fnv"
	"fmt"
	"math/big"
)

var hashFn = fnv.New64a()

func Hashing(s string, index int, size int64) uint {
	write := []byte(s + fmt.Sprint(index))
	hashFn.Write(write)
	usize := uint64(size)
	hash := hashFn.Sum64()
	hashFn.Reset()
	return uint(hash % usize % usize)
}

func (cb *cbBloomFilter) Check(index int, bf *big.Int) bool {
	val := bf.Bit(index)
	return val == 1
}
