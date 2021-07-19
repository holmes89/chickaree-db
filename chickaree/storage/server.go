package storage

import (
	"context"
	"sync"

	"github.com/holmes89/chickaree-db/chickaree"
	"github.com/holmes89/chickaree-db/chickaree/discovery"
	"github.com/rs/zerolog/log"
)

type Server struct {
	chickaree.UnimplementedChickareeDBServer
	ServerConfig
	store      *DistributedStorage
	membership *discovery.Membership

	shutdown     bool
	shutdowns    chan struct{}
	shutdownLock sync.Mutex
}

func NewServer(config ServerConfig) (*Server, error) {
	store, err := newStorage(config.StoragePath)
	if err != nil {
		return nil, err
	}
	distStore, err := NewDistributedStorage(store, config.Config)
	if err != nil {
		return nil, err
	}
	s := &Server{
		ServerConfig: config,
		store:        distStore,
	}

	s.setupMembership()

	return s, nil
}
func (s *Server) Close() error {
	log.Info().Msg("server closing...")
	s.shutdownLock.Lock()
	defer s.shutdownLock.Unlock()
	if s.shutdown {
		return nil
	}
	s.shutdown = true
	close(s.shutdowns)
	shutdown := []func() error{
		s.membership.Leave,
		s.store.Close,
	}

	for _, fn := range shutdown {
		if err := fn(); err != nil {
			return err
		}
	}
	log.Info().Msg("server closed.")
	return nil
}

func (s *Server) setupMembership() (err error) {
	log.Info().Msg("setting up membership")
	rpcAddr, err := s.ServerConfig.RPCAddr()
	if err != nil {
		return err
	}
	s.membership, err = discovery.New(s.store, discovery.Config{
		NodeName: s.ServerConfig.NodeName,
		BindAddr: s.ServerConfig.BindAddr,
		Tags: map[string]string{
			"rpc_addr": rpcAddr,
		},
		StartJoinAddrs: s.ServerConfig.StartJoinAddrs,
	})
	log.Info().Msg("membership established.")
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
	return &chickaree.SetResponse{}, err
}

func (s *Server) GetServers(
	ctx context.Context, req *chickaree.GetServersRequest,
) (
	*chickaree.GetServersResponse, error) {
	servers, err := s.store.GetServers()
	if err != nil {
		return nil, err
	}
	return &chickaree.GetServersResponse{Servers: servers}, nil
}
