package memtable

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
)

type LvlUp struct {
	ds List
}

func (l *LvlUp) Put(key, value []byte) error {
	// TODO what are all of the conditions that would cause an error?
	if l.ds.Header == nil {
		return errors.New("no data store available")
	}
	l.ds.Insert(key, value)

	return nil
}

func (l *LvlUp) Has(key []byte) (bool, error) {
	// TODO what are all of the conditions that would cause an error?
	if l.ds.Header == nil {
		return false, errors.New("no data store available")
	}
	_, ok := l.ds.Search(key)

	return ok, nil
}

func (l *LvlUp) Get(key []byte) ([]byte, error) {
	val, ok := l.ds.Search(key)
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf("cannot get %s", string(key))
}

func (l *LvlUp) Delete(key []byte) error {
	if ok, _ := l.Has(key); ok {
		l.ds.Delete(key)
		return nil
	}
	return fmt.Errorf("%s not present", string(key))
}

func (l *LvlUp) RangeScan(start, limit []byte) (Iterator, error) {
	i := &Iter{map[string][]byte{}, []string{}, 0}

	// TODO this is not what I want to actually range over;
	//		it misses the rest of the list
	for _, n := range l.ds.Header.Forward {
		k := n.Key
		fmt.Println(k)
		// start: inclusive
		// limit: exclusive
		if bytes.Compare(start, k) <= 0 && bytes.Compare(k, limit) == -1 {
			i.Pairs[string(k)] = n.Value
			i.Keys = append(i.Keys, string(k))
		}
	}

	sort.Strings(i.Keys)

	return i, nil
}

type Iter struct {
	Pairs map[string][]byte
	Keys []string
	Index  int
}

func (i *Iter) Next() bool {
	keysLen := len(i.Keys)

	if i.Index == keysLen {
		return false
	}

	i.Index += 1
	return true
}

// Error returns any accumulated error. Exhausting all the key/value pairs
// is not considered to be an error.
func (i *Iter) Error() error {
	// TODO
	return nil
}

func (i *Iter) Key() []byte {
	ind := i.Index
	keys := i.Keys

	if ind >= len(keys) {
		return nil
	}

	return []byte(keys[ind])
}

func (i *Iter) Value() []byte {
	ind := i.Index
	keys := i.Keys

	fmt.Println("ind", ind, "keys", keys)

	if ind >= len(i.Pairs) {
		return nil
	}

	return i.Pairs[keys[ind]]
}
