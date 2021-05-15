package core

import (
	"errors"
	"log"
	"strings"

	"github.com/boltdb/bolt"
)

var (
	ErrNotFound = errors.New("entry not found")
)

type Entry struct {
	Type  string
	Key   []byte
	Value []byte
	Args  [][]byte
}

type EntryType byte

const (
	unknownEntry EntryType = iota
	primativeEntry
	mapEntry
	arrayEntry
	setEntry
)

func encodeEntry(e Entry) []byte {
	var t []byte
	switch strings.ToLower(e.Type) {
	case "primative":
		t = append(t, byte(primativeEntry))
	case "map":
		t = append(t, byte(mapEntry))
	case "array":
		t = append(t, byte(arrayEntry))
	case "set":
		t = append(t, byte(setEntry))
	default:
		t = append(t, byte(unknownEntry))
	}
	return append(t, e.Value...)
}

func decodeEntry(key, value []byte) Entry {
	e := Entry{
		Key:   key,
		Value: value[1:],
	}
	switch EntryType(value[0]) {
	case primativeEntry:
		e.Type = "primative"
	case mapEntry:
		e.Type = "map"
	case arrayEntry:
		e.Type = "array"
	case setEntry:
		e.Type = "set"
	default:
		e.Type = "unknown"
	}
	return e
}

type Repository interface {
	Set(Entry) error
	Get([]byte) (Entry, error) // future  args ...[]byte for faster searching in list or map
	All() <-chan Entry
	Remove([]byte) error
	Close() error
}

type repo struct {
	db *bolt.DB
}

var defaultBucket = []byte{0x0}

func NewRepo(f string) Repository {
	db, err := bolt.Open(f, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(defaultBucket)
		return err
	})
	return &repo{
		db: db,
	}
}

func (r *repo) Close() error {
	return r.db.Close()
}

func (r *repo) Set(e Entry) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		v := encodeEntry(e)
		return tx.Bucket(defaultBucket).Put(e.Key, v)
	})
}

func (r *repo) Get(key []byte) (Entry, error) {
	var b []byte
	r.db.View(func(tx *bolt.Tx) error {
		b = tx.Bucket(defaultBucket).Get(key)
		return nil
	})

	if len(b) == 0 {
		return Entry{}, ErrNotFound
	}
	return decodeEntry(key, b), nil
}

func (r *repo) All() <-chan Entry {
	ch := make(chan Entry)
	go r.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(defaultBucket).ForEach(func(k, v []byte) error {
			ch <- decodeEntry(k, v)
			return nil
		})
	})
	return ch
}

func (r *repo) Remove(key []byte) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)

		if b.Get(key) == nil {
			return nil
		}

		if err := b.Delete(key); err != nil {
			return err
		}
		return nil
	})
}
