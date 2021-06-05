package storage

import (
	"context"

	"github.com/holmes89/chickaree-db/chickaree"
)

type Config struct {
}

type server struct {
	chickaree.UnimplementedChickareeDBServer
	*Config
	store storage
}

type Closer interface {
	Close() error
}

func NewServer(config *Config) (srv chickaree.ChickareeDBServer, closer Closer, err error) {
	store, err := newStorage("chickaree.db")
	if err != nil {
		return nil, nil, err
	}
	s := &server{
		Config: config,
		store:  store,
	}
	return s, s, nil
}

func (s *server) Close() error {
	return s.store.Close()
}

func (s *server) Get(ctx context.Context, req *chickaree.GetRequest) (*chickaree.GetResponse, error) {
	v, err := s.store.Get([]byte(req.Key))

	if err != nil {
		return nil, err
	}

	resp := &chickaree.GetResponse{
		Data: v,
	}
	return resp, nil
}

func (s *server) Set(ctx context.Context, req *chickaree.SetRequest) (*chickaree.SetResponse, error) {
	err := s.store.Set([]byte(req.Key), req.Value)

	if err != nil {
		return nil, err
	}

	return &chickaree.SetResponse{}, nil
}
