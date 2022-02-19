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
	fwd := BuildForwardList(1)
	header := &Node{fwd, nil, nil}
	l := List{header, 1}
	randLev := 6

	l.Insert([]byte("abc"), []byte("123"))
	l.Insert([]byte("def"), []byte("456"))

	if l.Level != randLev {
		t.Fatalf("Should have been %d, but was %d", randLev, l.Level)
	}

	val, ok := l.Search([]byte("def"))
	if !ok {
		t.Fatalf("Did not retrieve %s, but instead %s.", []byte("456"), val)
	}
	if !bytes.Equal([]byte("456"), val) {
		t.Fatalf("wtf? %s", val)
	}

	// test updating
	l.Insert([]byte("def"), []byte("789"))

	val, ok = l.Search([]byte("def"))
	if !ok {
		t.Fatalf("Did not retrieve %s, but instead %s.", []byte("789"), val)
	}
	if !bytes.Equal([]byte("789"), val) {
		t.Fatalf("wtf? %s", val)
	}
}

func TestSearch(t *testing.T) {
	// "The header of a list has forward pointers at levels one through MaxLevel."
	// But also, under the Init section, it says start with 1.
	headerForward := BuildForwardList(1)
	if len(headerForward) != 1 {
		t.Fatalf("Expected %d nodes in the forward list, but got %d.\n", MaxLevel, len(headerForward))
	}
	header := Node{headerForward, nil, nil}
	list := List{&header, 1}

	insertVal0 := []byte("345")
	list.Insert([]byte("def"), insertVal0)

	insertVal1 := []byte("678")
	list.Insert([]byte("ghi"), insertVal1)

	insertVal2 := []byte("901")
	list.Insert([]byte("jkl"), insertVal2)

	val, ok := list.Search([]byte("def"))
	if !ok {
		t.Fatalf("Did not retrieve %s, but instead %s.", insertVal0, val)
	}
	if !bytes.Equal(insertVal0, val) {
		t.Fatalf("wtf? %s", val)
	}
}

func TestDelete(t *testing.T) {
	fwd := BuildForwardList(1)
	header := &Node{fwd, nil, nil}
	l := List{header, 1}

	l.Insert([]byte("abc"), []byte("123"))
	l.Insert([]byte("def"), []byte("456"))
	l.Insert([]byte("ghi"), []byte("789"))

	l.Delete([]byte("def"))

	val, ok := l.Search([]byte("def"))
	if ok {
		t.Fatalf("We got back %c, when we were expecting nada", val)
	}

}
