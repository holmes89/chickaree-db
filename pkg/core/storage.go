package core

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

type Arg []byte
type Repository interface {
	Set(key []Arg) Response
	Get(key []Arg) Response
	Del(key []Arg) Response
	HSet(key []Arg) Response
	HGet(key []Arg) Response
	HGetAll(key []Arg) ResponseArray
	HExists(key []Arg) Response
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

func (r *repo) HSet(args []Arg) Response {
	var count int
	if (len(args) < 3) || (len(args[1:])%2 != 0) {
		return ErrResponse(fmt.Errorf("invalid arg count %d", len(args)))
	}
	err := r.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(args[0])
		if err != nil {
			return err
		}
		for i := 1; i < len(args); i += 2 {
			if err := b.Put(args[i], args[i+1]); err != nil {
				return err
			}
			count++
		}
		return nil
	})
	if err != nil {
		ErrResponse(err)
	}
	c := fmt.Sprintf("%d", count)
	return Response{
		rtype:   Integers,
		content: []byte(c),
	}
}

func (r *repo) HGet(args []Arg) Response {
	var b []byte
	r.db.View(func(tx *bolt.Tx) error {
		b = tx.Bucket(args[0]).Get(args[1])
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

// TODO maybe need to store types along with value
func (r *repo) HGetAll(args []Arg) ResponseArray {
	var res ResponseArray
	r.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(args[0]).ForEach(func(k, v []byte) error {
			res = append(res, Response{
				rtype:   BulkStrings,
				content: k,
				length:  len(k),
			})
			res = append(res, Response{
				rtype:   BulkStrings,
				content: v,
				length:  len(v),
			})
			return nil
		})
	})

	return res
}

func (r *repo) HExists(args []Arg) Response {

	var b []byte
	r.db.View(func(tx *bolt.Tx) error {
		b = tx.Bucket(args[0]).Get(args[1])
		return nil
	})
	count := "1"
	if len(b) == 0 {
		count = "0"
	}
	res := Response{
		rtype:   Integers,
		content: []byte(count),
	}
	return res
}

func (r *repo) Del(args []Arg) Response {
	var count int
	err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)

		for _, a := range args {
			if b.Get(a) == nil {
				continue
			}
			if err := b.Delete(a); err != nil {
				return err
			}
			count++
		}
		return nil
	})
	if err != nil {
		ErrResponse(err)
	}
	c := fmt.Sprintf("%d", count)
	return Response{
		rtype:   Integers,
		content: []byte(c),
	}
}
