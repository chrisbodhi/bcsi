package memtable

import (
	"bytes"
	"fmt"
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
	// Level	int // the paper says we don't need to store the level of the node in the node (p. 1)
}

func (list *List) Search(searchKey []byte) ([]byte, bool) {
	x := list.Header
	for i := list.Level; i >= 0; i-- {
		for bytes.Compare(x.Forward[i].Key, searchKey) == -1 {
			x = x.Forward[i]
		}
	}

	x = x.Forward[0]
	if bytes.Equal(x.Key, searchKey) {
		return x.Value, true
	}

	return []byte{}, false
}

func (list *List) Insert(searchKey, newValue []byte) {
	var update [MaxLevel](*Node)
	fmt.Printf("Search key at start of insert: %s\n", searchKey)

	x := list.Header
	for i := list.Level; i >= 0; i-- {
		for len(x.Forward) > 0 && len(x.Forward[i].Key) > 0 && bytes.Compare(x.Forward[i].Key, searchKey) == -1 {
			x = x.Forward[i]
		}
		update[i] = x
	}

	x = x.Forward[0]
	if bytes.Equal(x.Key, searchKey) {
		x.Value = newValue
	} else {
		lvl := randomLevel()
		if lvl > list.Level {
			for i := list.Level + 1; i <= lvl; i++ {
				update[i] = list.Header
			}
			list.Level = lvl
		}
		x = makeNode(lvl, searchKey, newValue) // newValue is written just as "value" in the paper -- this is a guess
		for i := 0; i < list.Level; i++ {      // list.Level is written just as "level" in the paper -- this is a guess
			fmt.Println("len x fwd:", len(x.Forward), "update:", update[i].Forward)
			x.Forward[i] = update[i].Forward[i]
			update[i].Forward[i] = x
		}
	}
	// TODO how does this do its update in-place?
}

// func (list *List) Delete(searchKey []byte) {
// 	var update [MaxLevel](*Node)

// 	x := list.Header
// 	for i := list.Level; i >= 0; i-- {
// 		for bytes.Compare(x.Forward[i].Key, searchKey) == -1 {
// 			x = x.Forward[i]
// 		}
// 		update[i] = x
// 	}
// 	x = x.Forward[1]
// 	if bytes.Equal(x.Key, searchKey) {
// 		for i := 1; i <= list.Level; i++ {
// 			if !areEqual(update[i].Forward[i], x) {
// 				break
// 			}
// 			update[i].Forward[i] = x.Forward[i]
// 		}
// 		for list.Level > 0 && isEmpty(list.Header.Forward[list.Level]) {
// 			list.Level -= 1
// 		}
// 	}
// }

func randomLevel() int {
	lvl := 0
	for rand.Float32() < p && lvl < MaxLevel {
		lvl += 1
	}
	return lvl
}

func makeNode(lvl int, searchKey, value []byte) *Node {
	fwdList := BuildForwardList(lvl)
	return &Node{Forward: fwdList, Key: searchKey, Value: value}
}

func areEqual(a, b Node) bool {
	return bytes.Equal(a.Key, b.Key) && bytes.Equal(a.Value, b.Value)
}

func isEmpty(n Node) bool {
	return len(n.Forward) == 0 && len(n.Key) == 0 && len(n.Value) == 0
}

func BuildForwardList(lvl int) [](*Node) {
	fwdList := [](*Node){&Node{}}

	for i := 1; i < lvl; i++ {
		fwdList = append(fwdList, &Node{})
	}
	fmt.Println("lvl", lvl, "fwdList", len(fwdList))

	return fwdList
}
