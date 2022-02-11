package memtable

import (
	"testing"
)

func TestGet(t *testing.T) {
	var l LvlUp
	l.ds = make(map[string][]byte)

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

func TestRangeScan(t *testing.T) {
	var l LvlUp
	l.ds = make(map[string][]byte)

	for i := 'A'; i <= 'Z'; i++ {
		l.ds[string(i)] = []byte{byte(i + 5)}
	}


	start := []byte("B") // inclusive
	limit := []byte("E") // exclusive
	iter, err := l.RangeScan(start, limit)

	if err != nil {
		t.Fatalf("RangeScan err: %s\n", err)
	}

	firstValue := iter.Value()

	if firstValue[0] != start[0] {
		t.Fatalf("First value was not %c, it was %c\n", firstValue[0], start[0])
	}

	iter.Next() // 'C': 'C' + 5 = 'H'
	iter.Next() // 'D': 'D' + 5 = 'I'

	lastValue := iter.Value()

	// Minus 1 because exclusive
	if lastValue[0] != limit[0]-1 {
		t.Fatalf("Last key was not %c, it was %c\n", lastValue[0], limit[0]-1)
	}
}
