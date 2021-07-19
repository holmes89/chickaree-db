package storage

import (
	bolt "go.etcd.io/bbolt"
)

type storage interface {
	Set(key, value []byte) error
	Get(key []byte) ([]byte, error)
	Close() error
}

var defaultBucket = []byte{0x0}

type store struct {
	db   *bolt.DB
	path string
}

func newStorage(path string) (storage, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(defaultBucket)
		return err
	}); err != nil {
		return nil, err
	}

	return &store{
		path: path,
		db:   db,
	}, nil
}

func (s *store) Close() error {
	return s.db.Close()
}

func (s *store) Set(key, value []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(defaultBucket).Put(key, value)
	})
}
func (s *store) Get(key []byte) (res []byte, err error) {
	s.db.View(func(tx *bolt.Tx) error {
		res = tx.Bucket(defaultBucket).Get(key)
		return nil
	})
	return res, nil
}
