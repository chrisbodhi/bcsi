package memtable

import (
	"bytes"
	"math/rand"
)

const p = 0.25 // pg. 5
const MaxLevel = 8

// "Because these data structures are linked lists with extra pointers that skip over intermediate nodes, I named them _skip lists_."

type List struct {
	Header *Node
	Level  int
}

type Node struct {
	Forward [](*Node) // len(Forward) is this node's level: "A node that has k forward pointers is called a level k node."
	Key     []byte
	Value   []byte
	// Level	int	  // the paper says we don't need to store the level of the node in the node (p. 1)
}

func (list *List) Search(searchKey []byte) ([]byte, bool) {
	x := list.Header
	for i := (list.Level - 1); i >= 0; i-- {
		for i < len(x.Forward) && bytes.Compare(x.Forward[i].Key, searchKey) == -1 {
			x = x.Forward[i]
		}
	}

	x = x.Forward[0] // the base level
	if bytes.Equal(x.Key, searchKey) {
		return x.Value, true
	}

	return nil, false
}

func (list *List) Insert(searchKey, newValue []byte) {
	var update [MaxLevel](*Node)

	x := list.Header

	// List level of 1 means that a forward list has a length of 1 (and only an index of 0)
	for i := (list.Level - 1); i >= 0; i-- {
		for len(x.Forward) > 0 &&
			len(x.Forward[i].Key) > 0 &&
			bytes.Compare(x.Forward[i].Key, searchKey) == -1 {
			x = x.Forward[i]
		}
		update[i] = x
	}

	x = x.Forward[0]

	if bytes.Equal(x.Key, searchKey) {
		// Update the value for searchKey; it's already in the skip list
		x.Value = newValue
	} else {
		lvl := randomLevel()
		if lvl > list.Level {
			for i := list.Level; i < lvl; i++ {
				update[i] = list.Header
				list.Header.Forward = append(list.Header.Forward, makeNode(i, nil, nil))
			}
			list.Level = lvl
		}
		x = makeNode(list.Level, searchKey, newValue) // newValue is written just as "value" in the paper -- this is a guess
		for i := 0; i < list.Level; i++ {       // initLevel is written just as "level" in the paper -- this is a guess
			x.Forward[i] = update[i].Forward[i]
			update[i].Forward[i] = x
		}
	}
}

func (list *List) Delete(searchKey []byte) {
	var update [MaxLevel](*Node)

	x := list.Header

	for i := (list.Level - 1); i >= 0; i-- {
		for bytes.Compare(x.Forward[i].Key, searchKey) == -1 {
			x = x.Forward[i]
		}
		update[i] = x
	}

	x = x.Forward[0]

	if bytes.Equal(x.Key, searchKey) {
		for i := 0; i < list.Level; i++ {
			if !areEqual(update[i].Forward[i], x) {
				break
			}
			update[i].Forward[i] = x.Forward[i]
		}
		for list.Level > 0 && isEmpty(list.Header.Forward[list.Level - 1]) {
			list.Level -= 1
		}
	}
}

func randomLevel() int {
	lvl := 0
	for rand.Float32() > p && lvl < MaxLevel {
		lvl += 1
	}
	return lvl
}

func makeNode(lvl int, searchKey, value []byte) *Node {
	fwdList := BuildForwardList(lvl)
	return &Node{Forward: fwdList, Key: searchKey, Value: value}
}

func areEqual(a, b *Node) bool {
	return bytes.Equal(a.Key, b.Key) && bytes.Equal(a.Value, b.Value)
}

func isEmpty(n *Node) bool {
	return len(n.Forward) == 0 && len(n.Key) == 0 && len(n.Value) == 0
}

func BuildForwardList(lvl int) [](*Node) {
	fwdList := [](*Node){&Node{}}

	for i := 1; i < lvl; i++ {
		fwdList = append(fwdList, &Node{})
	}

	return fwdList
}
