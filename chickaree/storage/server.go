package storage

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"sync"

	"github.com/holmes89/chickaree-db/chickaree"
	"github.com/holmes89/chickaree-db/chickaree/discovery"
)

type ServerConfig struct {
	Config
	ServerTLSConfig *tls.Config
	PeerTLSConfig   *tls.Config
	// DataDir stores the log and raft data.
	DataDir string
	// BindAddr is the address serf runs on.
	BindAddr string
	// RPCPort is the port for client (and Raft) connections.
	RPCPort int
	// Raft server id.
	NodeName string
	// Bootstrap should be set to true when starting the first node of the cluster.
	StartJoinAddrs []string
	ACLModelFile   string
	ACLPolicyFile  string
	Bootstrap      bool
}

func (c ServerConfig) RPCAddr() (string, error) {
	host, _, err := net.SplitHostPort(c.BindAddr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", host, c.RPCPort), nil
}

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
	return nil
}

func (s *Server) setupMembership() (err error) {
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
