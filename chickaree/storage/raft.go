package storage

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
)

var _ raft.StreamLayer = &StreamLayer{}

func NewRaft(ctx context.Context, cfg Config, myID, myAddress string, fsm raft.FSM) (*raft.Raft, *raft.NetworkTransport, error) {
	c := raft.DefaultConfig()
	c.LocalID = cfg.Raft.LocalID

	config := raft.DefaultConfig()
	config.LocalID = cfg.Raft.LocalID
	if cfg.Raft.HeartbeatTimeout != 0 {
		config.HeartbeatTimeout = cfg.Raft.HeartbeatTimeout
	}
	if cfg.Raft.ElectionTimeout != 0 {
		config.ElectionTimeout = cfg.Raft.ElectionTimeout
	}
	if cfg.Raft.LeaderLeaseTimeout != 0 {
		config.LeaderLeaseTimeout = cfg.Raft.LeaderLeaseTimeout
	}
	if cfg.Raft.CommitTimeout != 0 {
		config.CommitTimeout = cfg.Raft.CommitTimeout
	}

	baseDir := filepath.Join(cfg.RaftDir, myID)

	ldb, err := boltdb.NewBoltStore(filepath.Join(baseDir, "logs.dat"))
	if err != nil {
		return nil, nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "logs.dat"), err)
	}

	sdb, err := boltdb.NewBoltStore(filepath.Join(baseDir, "stable.dat"))
	if err != nil {
		return nil, nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "stable.dat"), err)
	}

	fss, err := raft.NewFileSnapshotStore(baseDir, 3, os.Stderr)
	if err != nil {
		return nil, nil, fmt.Errorf(`raft.NewFileSnapshotStore(%q, ...): %v`, baseDir, err)
	}

	maxPool := 5
	timeout := 10 * time.Second
	tm := raft.NewNetworkTransport(cfg.Raft.StreamLayer, maxPool, timeout, os.Stderr)

	r, err := raft.NewRaft(c, fsm, ldb, sdb, fss, tm)
	if err != nil {
		return nil, nil, fmt.Errorf("raft.NewRaft: %v", err)
	}

	hasState, err := raft.HasExistingState(
		ldb,
		sdb,
		fss,
	)
	if err != nil {
		return r, tm, err
	}
	if cfg.Raft.Bootstrap && !hasState {
		config := raft.Configuration{
			Servers: []raft.Server{{
				ID:      config.LocalID,
				Address: raft.ServerAddress(cfg.Raft.BindAddr),
			}},
		}
		err = r.BootstrapCluster(config).Error()
	}

	return r, tm, nil
}

type StreamLayer struct {
	ln              net.Listener
	serverTLSConfig *tls.Config
	peerTLSConfig   *tls.Config
}

func NewStreamLayer(
	ln net.Listener,
	serverTLSConfig,
	peerTLSConfig *tls.Config,
) *StreamLayer {
	return &StreamLayer{
		ln:              ln,
		serverTLSConfig: serverTLSConfig,
		peerTLSConfig:   peerTLSConfig,
	}
}

const RaftRPC = 1

func (s *StreamLayer) Dial(
	addr raft.ServerAddress,
	timeout time.Duration,
) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: timeout}
	var conn, err = dialer.Dial("tcp", string(addr))
	if err != nil {
		return nil, err
	}
	// identify to mux this is a raft rpc
	_, err = conn.Write([]byte{byte(RaftRPC)})
	if err != nil {
		return nil, err
	}
	if s.peerTLSConfig != nil {
		conn = tls.Client(conn, s.peerTLSConfig)
	}
	return conn, err
}

func (s *StreamLayer) Accept() (net.Conn, error) {
	conn, err := s.ln.Accept()
	if err != nil {
		return nil, err
	}
	b := make([]byte, 1)
	_, err = conn.Read(b)
	if err != nil {
		return nil, err
	}
	if bytes.Compare([]byte{byte(RaftRPC)}, b) != 0 {
		return nil, fmt.Errorf("not a raft rpc")
	}
	if s.serverTLSConfig != nil {
		return tls.Server(conn, s.serverTLSConfig), nil
	}
	return conn, nil
}

func (s *StreamLayer) Close() error {
	return s.ln.Close()
}

func (s *StreamLayer) Addr() net.Addr {
	return s.ln.Addr()
}
