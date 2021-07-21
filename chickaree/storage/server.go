package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	"github.com/holmes89/chickaree-db/chickaree"
	"github.com/holmes89/chickaree-db/chickaree/discovery"
	"github.com/rs/zerolog/log"
	"github.com/soheilhy/cmux"
)

type Server struct {
	chickaree.UnimplementedChickareeDBServer
	ServerConfig
	store      *DistributedStorage
	membership *discovery.Membership
	mux        cmux.CMux

	shutdown     bool
	shutdowns    chan struct{}
	shutdownLock sync.Mutex
}

func NewServer(config ServerConfig) (*Server, error) {

	s := &Server{
		ServerConfig: config,
		shutdowns:    make(chan struct{}),
	}

	setup := []func() error{
		s.setupMux,
		s.setupStorage,
		s.setupMembership,
	}
	for _, fn := range setup {
		if err := fn(); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Server) Mux() net.Listener {
	return s.mux.Match(cmux.Any())
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
	s.mux.Close()
	log.Info().Msg("server closed.")
	return nil
}

func (s *Server) setupMembership() (err error) {
	log.Info().Str("bind-addr", s.ServerConfig.BindAddr).Str("node-name", s.ServerConfig.NodeName).Msg("setting up membership")
	rpcAddr, err := s.ServerConfig.RPCAddr()
	if err != nil {
		log.Error().Err(err).Msg("unable to resolve rpc")
		return errors.New("unable to setup membership")
	}
	s.membership, err = discovery.New(s.store, discovery.Config{
		NodeName: s.ServerConfig.NodeName,
		BindAddr: s.ServerConfig.BindAddr,
		Tags: map[string]string{
			"rpc_addr": rpcAddr,
		},
		StartJoinAddrs: s.ServerConfig.StartJoinAddrs,
	})
	if err != nil {
		log.Error().Err(err).Msg("unable to setup membership")
		return errors.New("unable to setup membership")
	}
	log.Info().Msg("membership established.")

	return nil
}

func (s *Server) setupStorage() error {
	store, err := newStorage(s.ServerConfig.StoragePath)
	if err != nil {
		log.Error().Err(err).Msg("unable to create storage")
		return errors.New("unable to setup storage")
	}
	raftLn := s.mux.Match(func(reader io.Reader) bool {
		b := make([]byte, 1)
		if _, err := reader.Read(b); err != nil {
			return false
		}
		return bytes.Compare(b, []byte{byte(RaftRPC)}) == 0
	})
	config := s.ServerConfig.Config
	config.Raft.StreamLayer = NewStreamLayer(
		raftLn,
		s.ServerConfig.ServerTLSConfig,
		s.ServerConfig.PeerTLSConfig,
	)
	rpcAddr, err := s.ServerConfig.RPCAddr()
	if err != nil {
		log.Error().Err(err).Msg("unable to get rpc addr")
		return errors.New("unable to setup storage")
	}
	config.Raft.BindAddr = rpcAddr
	config.Raft.LocalID = raft.ServerID(s.ServerConfig.NodeName)
	config.Raft.Bootstrap = s.ServerConfig.Bootstrap
	s.store, err = NewDistributedStorage(
		store,
		config,
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to create distributed storage")
		return errors.New("unable to setup storage")
	}
	if s.ServerConfig.Bootstrap {
		if err := s.store.WaitForLeader(3 * time.Second); err != nil {
			log.Error().Err(err).Msg("unable to wait for leader")
			return errors.New("unable to setup storage")
		}
	}
	return nil
}

func (s *Server) setupMux() error {
	log.Info().Str("bind-addr", s.ServerConfig.BindAddr).Msg("creating mux...")
	addr, err := net.ResolveTCPAddr("tcp", s.ServerConfig.BindAddr)
	if err != nil {
		log.Error().Err(err).Str("bind-addr", s.ServerConfig.BindAddr).Msg("failed to setup tcp")
		return errors.New("unable to setup mux")
	}
	rpcAddr := fmt.Sprintf(
		"%s:%d",
		addr.IP.String(),
		s.ServerConfig.RPCPort,
	)
	ln, err := net.Listen("tcp", rpcAddr)
	if err != nil {
		log.Error().Err(err).Str("rpc-addr", rpcAddr).Msg("failed to listen")
		return errors.New("unable to setup mux")
	}
	s.mux = cmux.New(ln)
	log.Info().Msg("mux created.")
	return nil
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
