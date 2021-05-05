package core

import (
	"log"

	"github.com/boltdb/bolt"
)

type Arg []byte
type Repository interface {
	Set(key []Arg) Response
	Get(key []Arg) Response
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

func (r *repo) Set(args []Arg) Response {
	err := r.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(defaultBucket).Put(args[0], args[1])
	})
	if err != nil {
		return ErrResponse(err)
	}
	return OkResp
}

func (r *repo) Get(args []Arg) Response {
	var b []byte
	r.db.View(func(tx *bolt.Tx) error {
		b = tx.Bucket(defaultBucket).Get(args[0])
		return nil
	})
	if len(b) == 0 {
		return NilStringResp
	}
	res := Response{
		rtype:   BulkStrings,
		content: b,
		length:  len(b),
	}
	return res
}
