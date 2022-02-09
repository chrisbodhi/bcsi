package memtable

import "testing"

func TestGet(t *testing.T) {
	var l LvlUp
	l.ds = make(map[string][]byte)

	k := []byte("do")
	v := []byte("dope")

	l.Put(k, v)
	do, err := l.Get(k)

	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	for i, b := range do {
		if b != v[i] {
			t.Fatalf("Got %c, expected %c", b, v[i])
		}
	}
}