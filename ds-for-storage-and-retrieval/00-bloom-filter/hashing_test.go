package main

import (
	"testing"
)

var cb cbBloomFilter

func TestHashing(t *testing.T) {
	// 0.1 gets us two hashFns
	// 14 gets us size 67
	cb = *newCbBloomFilter(0.1, 14)
	h1 := Hashing("aardwolf", 0, int(cb.data.Int64()))
	h2 := Hashing("aardwolf", 1, int(cb.data.Int64()))

	if h1 == h2 {
		t.Errorf("Should not be equal: %d %d\n", h1, h2)
	}
}
