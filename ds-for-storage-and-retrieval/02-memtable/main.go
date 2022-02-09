package main

import (
	"errors"
	"fmt"
	"log"
)

type DB interface {
	// Get gets the value for the given key. It returns an error if the
	// DB does not contain the key.
	Get(key []byte) (value []byte, err error)

	// Has returns true if the DB contains the given key.
	Has(key []byte) (ret bool, err error)

	// Put sets the value for the given key. It overwrites any previous value
	// for that key; a DB is not a multi-map.
	Put(key, value []byte) error

	// Delete deletes the value for the given key.
	Delete(key []byte) error

	// RangeScan returns an Iterator (see below) for scanning through all
	// key-value pairs in the given range, ordered by key ascending.
	RangeScan(start, limit []byte) (Iterator, error)
}

type LvlUp struct {
	ds map[string][]byte
}

func (l *LvlUp) Put(key, value []byte) error {
	// TODO what are the conditions that would cause an error?
	if l.ds == nil {
		return errors.New("no data store available")
	}
	l.ds[string(key)] = value

	return nil
}

func (l *LvlUp) Has(key []byte) (bool, error) {
	// TODO what are the conditions that would cause an error?
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
		// TODO actually delete
		return nil
	} else {
		return fmt.Errorf("%s not present", string(key))
	}
}

type Iterator interface {
	// Next moves the iterator to the next key/value pair.
	// It returns false if the iterator is exhausted.
	Next() bool

	// Error returns any accumulated error. Exhausting all the key/value pairs
	// is not considered to be an error.
	Error() error

	// Key returns the key of the current key/value pair, or nil if done.
	Key() []byte

	// Value returns the value of the current key/value pair, or nil if done.
	Value() []byte
}

func main() {
	var l LvlUp
	l.ds = make(map[string][]byte)

	k := []byte{'d', 'o'}
	v := []byte{'d', 'o', 'p', 'e'}

	l.Put(k, v)
	do, err := l.Get(k)

	if err != nil {
		log.Fatal("dammit", err)
	}

	fmt.Println("do", string(do), "fuck hope")
}
