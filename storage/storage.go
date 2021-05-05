package storage

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

type Arg []byte
type Repository interface {
	Set(args []Arg) error
	Get(args []Arg) ([]byte, error)
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

func (r *repo) Set(args []Arg) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(defaultBucket).Put(args[0], args[1])
	})
}

func (r *repo) Get(args []Arg) (b []byte, err error) {
	r.db.View(func(tx *bolt.Tx) error {
		b = tx.Bucket(defaultBucket).Get(args[0])
		return nil
	})
	l := len(b)
	la := []byte(fmt.Sprintf("$%d\r\n", l))
	b = append(la, b...)
	b = append(b, []byte("\r\n")...)
	return b, nil
}
