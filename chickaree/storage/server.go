package storage

import (
	"context"
	"sync"

	"github.com/holmes89/chickaree-db/chickaree"
	"github.com/holmes89/chickaree-db/chickaree/discovery"
)

type Config struct {
}

type Server struct {
	chickaree.UnimplementedChickareeDBServer
	*Config
	store      storage
	events     *eventLogger
	membership *discovery.Membership // Refactor to extract replication logic.

	replicator *Replicator

	mu      sync.Mutex
	servers map[string]chan struct{}
	closed  bool
	close   chan struct{}
}

func NewServer(config *Config) (*Server, error) {
	store, err := newStorage("chickaree.db")
	if err != nil {
		return nil, err
	}
	s := &Server{
		Config: config,
		events: NewEventLogger(),
		store:  store,
	}

	return s, nil
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
	err := s.store.Set([]byte(req.Key), req.Value)
	if err != nil {
		return nil, err
	}
	go s.events.Write("set", [][]byte{req.Value})
	return &chickaree.SetResponse{}, nil
}

func (s *Server) EventLog(req *chickaree.EventLogRequest, stream chickaree.ChickareeDB_EventLogServer) error {
	evts, err := s.events.Events()
	if err != nil {
		return err
	}
	defer s.events.Leave(evts)
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case evt := <-evts:
			stream.Send(&chickaree.EventLogResponse{
				Command: []byte(evt.Command),
				Args:    evt.Args,
			})
		}
	}
}
