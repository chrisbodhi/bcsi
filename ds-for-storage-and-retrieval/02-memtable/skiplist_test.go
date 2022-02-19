package memtable

import (
	"bytes"
	"testing"
)

func TestBuildForwardList(t *testing.T) {
	levels := 0 // Position 0 is our first level
	list := BuildForwardList(levels)

	if len(list) != (levels+1) {
		t.Fatalf("Expected %d nodes in the forward list, but got %d.\n", (levels+1), len(list))
	}
}

func TestInsert(t *testing.T) {
	l := buildSkipList()

	// test updating
	l.Insert([]byte("def"), []byte("789"))

	val, ok := l.Search([]byte("def"))
	if !ok && !bytes.Equal([]byte("789"), val) {
		t.Fatalf("Did not retrieve %s, but instead %s.", []byte("789"), val)
	}
}

func TestSearch(t *testing.T) {
	l := buildSkipList()

	val, ok := l.Search([]byte("def"))

	if !ok && !bytes.Equal([]byte("456"), val) {
		t.Fatalf("Did not retrieve %s, but instead %s.", []byte("456"), val)
	}
}

func TestDelete(t *testing.T) {
	l := buildSkipList()

	l.Delete([]byte("def"))

	val, ok := l.Search([]byte("def"))
	if ok {
		t.Fatalf("We got back %c, when we were expecting nada", val)
	}

}

func buildSkipList() List {
	// "The header of a list has forward pointers at levels one through MaxLevel."
	// But also, under the Init section, it says start with 1.
	fwd := BuildForwardList(1)
	header := &Node{fwd, nil, nil}
	l := List{header, 1}

	l.Insert([]byte("abc"), []byte("123"))
	l.Insert([]byte("def"), []byte("456"))
	l.Insert([]byte("ghi"), []byte("789"))

	return l
}