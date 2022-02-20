package memtable

import (
	// "bytes"
	// "fmt"
	"log"
	"testing"
)

// Tests

func TestGet(t *testing.T) {
	var l LvlUp
	list, err := Init()
	if err != nil {
		log.Fatal(err)
	}

	l.ds = list

	k := []byte("do")
	v := []byte("dope")

	l.Put(k, v)
	do, err := l.Get(k)

	if err != nil {
		t.Fatalf("Get error: %s\n", err)
	}

	for i, b := range do {
		if b != v[i] {
			t.Fatalf("Got %c, expected %c", b, v[i])
		}
	}
}

func TestRangeScanSimple(t *testing.T) {
	var l LvlUp
	list, err := Init()
	if err != nil {
		log.Fatal(err)
	}

	l.ds = list

	for i := 'A'; i <= 'Z'; i++ {
		l.Put([]byte{byte(i)}, []byte{valFromKey(byte(i))})
	}

	start := []byte("B") // inclusive
	limit := []byte("E") // exclusive
	iter, err := l.RangeScan(start, limit)

	if err != nil {
		t.Fatalf("RangeScan err: %s\n", err)
	}

	firstValue := iter.Value()

	if firstValue[0] != valFromKey(start[0]) {
		t.Fatalf("First value was not %c, it was %c.\n", valFromKey(start[0]), firstValue[0])
	}

	iter.Next() // 'C': 'C' + 5 = 'H'
	iter.Next() // 'D': 'D' + 5 = 'I'

	lastValue := iter.Value()

	// Minus 1 because exclusive
	if lastValue[0] != valFromKey(limit[0]-1) {
		t.Fatalf("Last value was not %c, it was %c.\n", valFromKey(limit[0]-1), lastValue[0])
	}
}

func valFromKey(key byte) byte {
	return key + 5
}

// func TestRangeScanComplex(t *testing.T) {
// 	var l LvlUp
// 	l.ds = make(map[string][]byte)

// 	l.ds["A"] = []byte("b")
// 	l.ds["AA"] = []byte("after Z")
// 	l.ds["ABC"] = []byte("defg")
// 	l.ds["ABCD"] = []byte("apple")
// 	l.ds["B"] = []byte("boy")
// 	l.ds["BB"] = []byte("gun")
// 	l.ds["BFG"] = []byte("doom")
// 	l.ds["BG"] = []byte("bkgd")
// 	l.ds["C"] = []byte("d")
// 	l.ds["CAT"] = []byte("can")
// 	l.ds["CATCH"] = []byte("me")
// 	l.ds["E"] = []byte("h")
// 	l.ds["EEE"] = []byte("p")
// 	l.ds["EEK"] = []byte("the cat")

// 	start := []byte("BAT") // inclusive
// 	limit := []byte("DOG") // exclusive
// 	iter, err := l.RangeScan(start, limit)

// 	if err != nil {
// 		t.Fatalf("RangeScan err: %s\n", err)
// 	}

// 	firstValue := iter.Value()

// 	if !bytes.Equal(firstValue, l.ds["BB"]) {
// 		t.Fatalf("First value was not %c, it was %c.\n", l.ds["BB"], firstValue)
// 	}

// 	iter.Next()
// 	iter.Next()

// 	thirdValue := iter.Value()

// 	if !bytes.Equal(thirdValue, l.ds["BG"]) {
// 		t.Fatalf("Last value was not %c, it was %c.\n", l.ds["BG"], thirdValue)
// 	}

// 	// Let's run out the iterator

// 	if moreList := iter.Next(); !moreList {
// 		t.Fatalf("Should have kept going, instead cannot call Next: %t", moreList)
// 	}

// 	iter.Next()
// 	iter.Next()
// 	iter.Next()

// 	if moreToGo := iter.Next(); moreToGo {
// 		t.Fatalf("Should have ended, instead got %c | %t", iter.Key(), moreToGo)
// 	}

// 	if noVal := iter.Value(); noVal != nil {
// 		t.Fatalf("Should have failed, but we got %c", noVal)
// 	}
// }

// Benchmarks
// ----------
// For pre-Skip List implementation
// ➜ go test -run=XXX -bench .
// goos: darwin
// goarch: amd64
// pkg: github.com/chrisbodhi/bcsi/ds-for-storage-and-retrieval/memtable
// cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
// BenchmarkPut-4        	 1640446	       661.6 ns/op
// BenchmarkPutFixed-4   	42205840	        26.41 ns/op
// PASS
// ok  	github.com/chrisbodhi/bcsi/ds-for-storage-and-retrieval/memtable	4.302s

// fka BenchmarkPutFixed-4
func BenchmarkPut(b *testing.B) {
	var l LvlUp
	list, err := Init()
	if err != nil {
		log.Fatal(err)
	}
	l.ds = list
	k := []byte("A")
	v := []byte("AA")

	for i := 0; i < b.N; i++ {
		l.Put(k, v)
	}
}
