package storage

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	boltdb "github.com/hashicorp/raft-boltdb"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/raft"

	api "github.com/holmes89/chickaree-db/chickaree"
	"github.com/holmes89/chickaree-db/chickaree/discovery"
)

var logFile = "events.log"

type DistributedStorage struct {
	config Config
	store  storage
	raft   *raft.Raft
}

var (
	_ discovery.Handler = &DistributedStorage{}
)

func NewDistributedStorage(store storage, config Config) (*DistributedStorage, error) {
	l := &DistributedStorage{
		config: config,
		store:  store,
	}

	if err := l.setupRaft(); err != nil {
		return nil, err
	}
	return l, nil
}

func (s *DistributedStorage) setupRaft() error {
	baseDir := s.config.RaftDir

	logStore, err := boltdb.NewBoltStore(filepath.Join(baseDir, "logs.dat"))
	if err != nil {
		log.Error().Err(err).Msg("unable to create bolt sroage for logs")
		return fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "logs.dat"), err)
	}

	stableStore, err := boltdb.NewBoltStore(filepath.Join(baseDir, "stable.dat"))
	if err != nil {
		log.Error().Err(err).Msg("unable to create bolt sroage for stable messages")
		return fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "stable.dat"), err)
	}

	snapshotStore, err := raft.NewFileSnapshotStore(baseDir, 3, os.Stderr)
	if err != nil {
		log.Error().Err(err).Msg("unable to create bolt sroage for snapshots")
		return fmt.Errorf(`raft.NewFileSnapshotStore(%q, ...): %v`, baseDir, err)
	}

	log.Info().Msg("distributed server raft storage created")
	fsm := &fsm{store: s.store}

	maxPool := 5
	timeout := 30 * time.Second
	transport := raft.NewNetworkTransport(
		s.config.Raft.StreamLayer,
		maxPool,
		timeout,
		os.Stderr,
	)

	config := raft.DefaultConfig()
	config.LocalID = s.config.Raft.LocalID
	if s.config.Raft.HeartbeatTimeout != 0 {
		log.Info().Dur("timeout", s.config.Raft.HeartbeatTimeout).Msg("overriding heartbeat timeout")
		config.HeartbeatTimeout = s.config.Raft.HeartbeatTimeout
	}
	if s.config.Raft.ElectionTimeout != 0 {
		log.Info().Dur("timeout", s.config.Raft.ElectionTimeout).Msg("overriding election timeout")
		config.ElectionTimeout = s.config.Raft.ElectionTimeout
	}
	if s.config.Raft.LeaderLeaseTimeout != 0 {
		log.Info().Dur("timeout", s.config.Raft.LeaderLeaseTimeout).Msg("overriding leaderlease timeout")
		config.LeaderLeaseTimeout = s.config.Raft.LeaderLeaseTimeout
	}
	if s.config.Raft.CommitTimeout != 0 {
		log.Info().Dur("timeout", s.config.Raft.CommitTimeout).Msg("overriding commit timeout")
		config.CommitTimeout = s.config.Raft.CommitTimeout
	}

	s.raft, err = raft.NewRaft(
		config,
		fsm,
		logStore,
		stableStore,
		snapshotStore,
		transport,
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to create raft")
		return errors.New("failed to create raft")
	}
	hasState, err := raft.HasExistingState(
		logStore,
		stableStore,
		snapshotStore,
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine existing state")
		return errors.New("failed to create raft")
	}
	if s.config.Raft.Bootstrap && !hasState {
		log.Info().Interface("id", config.LocalID).Msg("bootstrapping cluster")
		config := raft.Configuration{
			Servers: []raft.Server{{
				ID:      config.LocalID,
				Address: raft.ServerAddress(s.config.Raft.BindAddr),
			}},
		}
		if err := s.raft.BootstrapCluster(config).Error(); err != nil {
			log.Error().Err(err).Msg("unable to bootstrap cluster")
			return err
		}
	}
	log.Info().Msg("distributed server raft connected")
	return nil
}

var _ storage = &DistributedStorage{}

type RequestType uint8

const (
	SetRequestType RequestType = 0
)

func (s *DistributedStorage) Set(key, value []byte) error {
	_, err := s.apply(SetRequestType, &api.SetRequest{
		Key:   string(key),
		Value: value,
	})
	if err != nil {
		return err
	}
	return s.store.Set(key, value)
}

func (s *DistributedStorage) Get(key []byte) ([]byte, error) {
	return s.store.Get(key)
}

func (s *DistributedStorage) apply(reqType RequestType, req proto.Message) (
	interface{},
	error,
) {
	var buf bytes.Buffer
	_, err := buf.Write([]byte{byte(reqType)})
	if err != nil {
		return nil, err
	}
	b, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	timeout := 30 * time.Second
	future := s.raft.Apply(buf.Bytes(), timeout)
	if future.Error() != nil {
		return nil, future.Error()
	}
	res := future.Response()
	if err, ok := res.(error); ok {
		return nil, err
	}
	return res, nil
}

func (s *DistributedStorage) Join(id, addr string) error {
	log.Info().Str("id", id).Str("addr", addr).Msg("joining cluster...")
	configFuture := s.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		log.Error().Err(err).Msg("unable to get raft configuration")
		return err
	}
	serverID := raft.ServerID(id)
	serverAddr := raft.ServerAddress(addr)
	for _, srv := range configFuture.Configuration().Servers {
		if srv.ID == serverID || srv.Address == serverAddr {
			if srv.ID == serverID && srv.Address == serverAddr {
				log.Warn().Str("id", id).Str("addr", addr).Msg("server already joined")
				return nil
			}
			// remove the existing server
			removeFuture := s.raft.RemoveServer(serverID, 0, 0)
			if err := removeFuture.Error(); err != nil {
				log.Error().Err(err).Msg("unable to remove existing server")
				return err
			}
		}
	}
	addFuture := s.raft.AddVoter(serverID, serverAddr, 0, 0)
	if err := addFuture.Error(); err != nil {
		log.Error().Err(err).Msg("unable to add voter")
		return err
	}
	log.Info().Str("id", id).Str("addr", addr).Msg("joined cluster")
	return nil
}

func (s *DistributedStorage) Leave(id string) error {
	removeFuture := s.raft.RemoveServer(raft.ServerID(id), 0, 0)
	return removeFuture.Error()
}

func (s *DistributedStorage) WaitForLeader(timeout time.Duration) error {
	timeoutc := time.After(timeout)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-timeoutc:
			log.Error().Msg("finding leader has timed out")
			return fmt.Errorf("timed out")
		case <-ticker.C:
			if l := s.raft.Leader(); l != "" {
				return nil
			}
		}
	}
}

func (s *DistributedStorage) Close() error {
	f := s.raft.Shutdown()
	if err := f.Error(); err != nil {
		log.Error().Err(err).Msg("unable to shutdown raft")
		return err
	}
	return s.store.Close()
}

func (s *DistributedStorage) GetServers(nameIP map[string]string) ([]*api.Server, error) {
	leaderIP := string(s.raft.Leader())

	future := s.raft.GetConfiguration()
	if err := future.Error(); err != nil {
		log.Error().Err(err).Msg("unable to get configuration")
		return nil, errors.New("unable to get servers")
	}

	var servers []*api.Server
	for _, server := range future.Configuration().Servers {
		ip := nameIP[string(server.Address)]
		log.Info().Str("addr", string(server.Address)).Str("ip", ip).Str("leader", leaderIP).Bool("is-leader", leaderIP == ip).Msg("server")
		servers = append(servers, &api.Server{
			Id:       string(server.ID),
			RpcAddr:  string(server.Address),
			IsLeader: leaderIP == ip,
		})
	}
	log.Info().Int("servers", len(servers)).Msg("returning servers")
	return servers, nil
}

var (
	_ raft.FSM = (*fsm)(nil)
)

type fsm struct {
	mu    sync.RWMutex
	store storage
}

func (s *fsm) Apply(record *raft.Log) interface{} {
	return s.apply(record.Data)
}

func (s *fsm) apply(buf []byte) interface{} {
	reqType := RequestType(buf[0])
	switch reqType {
	case SetRequestType:
		return s.applySet(buf[1:])
	}
	s.write(buf)
	return nil
}

func (s *fsm) applySet(b []byte) error {
	var req api.SetRequest
	err := proto.Unmarshal(b, &req)
	if err != nil {
		return err
	}

	if err := s.store.Set([]byte(req.Key), req.Value); err != nil {
		return err
	}
	return nil
}

func (s *fsm) write(body []byte) {
	body = append(body, '\n')
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open log file")
	}
	defer f.Close()

	if _, err := f.Write(body); err != nil {
		log.Fatal().Err(err).Msg("unable to write to log file")
	}
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	fi, err := os.OpenFile(logFile, os.O_RDONLY, 0600)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open log file")
		return nil, errors.New("failed to create snapshot")
	}

	defer fi.Close()
	b, err := io.ReadAll(fi)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to read log file")
		return nil, errors.New("failed to create snapshot")
	}
	return &snapshot{reader: bytes.NewBuffer(b)}, nil
}

func (f *fsm) Restore(r io.ReadCloser) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		f.apply(scanner.Bytes())
	}
	return nil
}

var _ raft.FSMSnapshot = (*snapshot)(nil)

type snapshot struct {
	reader io.Reader
}

func (s *snapshot) Persist(sink raft.SnapshotSink) error {
	if _, err := io.Copy(sink, s.reader); err != nil {
		_ = sink.Cancel()
		return err
	}
	return sink.Close()
}

func (s *snapshot) Release() {}

var _ raft.StreamLayer = (*StreamLayer)(nil)

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
		log.Error().Str("addr", string(addr)).Err(err).Msg("unable to dial")
		return nil, errors.New("failed to dial")
	}
	// identify to mux this is a raft rpc
	_, err = conn.Write([]byte{byte(RaftRPC)})
	if err != nil {
		log.Error().Err(err).Msg("unable to write raft rpc")
		return nil, errors.New("failed to dial")
	}
	if s.peerTLSConfig != nil {
		conn = tls.Client(conn, s.peerTLSConfig)
	}
	return conn, nil
}

func (s *StreamLayer) Accept() (net.Conn, error) {
	conn, err := s.ln.Accept()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to accept message")
		return nil, errors.New("failed to accept")
	}
	b := make([]byte, 1)
	_, err = conn.Read(b)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to read connection")
		return nil, errors.New("failed to accept")
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
