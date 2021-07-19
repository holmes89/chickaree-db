package storage

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	"github.com/holmes89/chickaree-db/chickaree"
	"github.com/holmes89/chickaree-db/chickaree/discovery"
	"google.golang.org/protobuf/proto"
)

type Config struct {
	StoragePath string
	RaftDir     string
	Raft        struct {
		raft.Config
		BindAddr    string
		StreamLayer *StreamLayer
		Bootstrap   bool
	}
}

type Server struct {
	chickaree.UnimplementedChickareeDBServer
	Config
	store      storage
	membership *discovery.Membership

	raft             *raft.Raft
	raftNetTransport *raft.NetworkTransport

	closeLock sync.Mutex
	servers   map[string]chan struct{}
	closed    bool
	close     chan struct{}
}

func NewServer(config Config) (*Server, error) {
	store, err := newStorage(config.StoragePath)
	if err != nil {
		return nil, err
	}
	s := &Server{
		Config: config,
		store:  store,
	}

	fsm := NewEventLogger(store)
	r, trans, err := NewRaft(context.Background(), config, "test", "test", fsm)
	if err != nil {
		return s, errors.New("unable to register raft protocol")
	}

	s.raft = r
	s.raftNetTransport = trans

	return s, nil
}
func (s *Server) Close() error {
	s.closeLock.Lock()
	defer s.closeLock.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	close(s.close)
	shutdown := []func() error{
		s.membership.Leave,
		s.store.Close,
	}

	for _, fn := range shutdown {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) encode(t RequestType, req proto.Message) error {
	var buf bytes.Buffer
	if err := buf.WriteByte(byte(t)); err != nil {
		return err
	}
	b, err := proto.Marshal(req)
	if err != nil {
		return err
	}
	if _, err := buf.Write(b); err != nil {
		return err
	}
	future := s.raft.Apply(buf.Bytes(), time.Second)
	if future.Error() != nil {
		return future.Error()
	}
	res := future.Response()
	if err, ok := res.(error); ok {
		return err
	}
	return nil
}

func (s *Server) setupMembership() (err error) {
	s.membership, err = discovery.New(s.replicator, discovery.Config{
		NodeName: "chickaree",
		BindAddr: "127.0.0.1:8081",
		Tags: map[string]string{
			"rpc_addr": "127.0.0.1:8080",
		},
		StartJoinAddrs: []string{},
	})
	return err
}

func (s *Server) Get(ctx context.Context, req *chickaree.GetRequest) (*chickaree.GetResponse, error) {
	v, err := s.store.Get([]byte(req.Key))

	if err != nil {
		return nil, err
	}

	resp := &chickaree.GetResponse{
		Data: v,
	}
	return resp, nil
}

func (s *Server) Set(ctx context.Context, req *chickaree.SetRequest) (*chickaree.SetResponse, error) {
	if err := s.encode(Set, req); err != nil {
		return nil, err
	}
	err := s.store.Set([]byte(req.Key), req.Value)
	if err != nil {
		return nil, err
	}
	return &chickaree.SetResponse{}, nil
}
