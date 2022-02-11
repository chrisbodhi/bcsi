package memtable

import (
	"errors"
	"fmt"
)

type LvlUp struct {
	ds map[string][]byte
}

func (l *LvlUp) Put(key, value []byte) error {
	// TODO what are all of the conditions that would cause an error?
	if l.ds == nil {
		return errors.New("no data store available")
	}
	l.ds[string(key)] = value

	return nil
}

func (l *LvlUp) Has(key []byte) (bool, error) {
	// TODO what are all of the conditions that would cause an error?
	if l.ds == nil {
		return false, errors.New("no data store available")
	}
	_, ok := l.ds[string(key)]

	return ok, nil
}

func (l *LvlUp) Get(key []byte) ([]byte, error) {
	val, ok := l.ds[string(key)]
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf("cannot get %s", string(key))
}

func (l *LvlUp) Delete(key []byte) error {
	if ok, _ := l.Has(key); ok {
		delete(l.ds, string(key))
		return nil
	}
	return fmt.Errorf("%s not present", string(key))
}

func (l *LvlUp) RangeScan(start, limit []byte) (Iterator, error) {
	i := &Iter{[][]byte{}, [][]byte{}, 0}
	// range over l.ds, starting from start, until
	// counter is >= current key
	// TODO write a helper to compare byte arrays
	return i, nil
}

type Iter struct {
	keys   [][]byte
	values [][]byte
	index  int
}

func (i *Iter) Next() bool {
	l := len(i.keys)
	nextIndex := i.index + 1

	if nextIndex >= l {
		return false
	}

	i.index = nextIndex
	return true
}

// Error returns any accumulated error. Exhausting all the key/value pairs
// is not considered to be an error.
func (i *Iter) Error() error {
	// TODO
	return nil
}

func (i *Iter) Key() []byte {
	ind := i.index
	keys := i.keys

	if len(keys) >= ind {
		return nil
	}

	return keys[ind]
}

func (i *Iter) Value() []byte {
	ind := i.index
	values := i.values

	if len(values) >= ind {
		return nil
	}

	return values[ind]
}
