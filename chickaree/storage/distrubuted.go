package storage

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

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
		return fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "logs.dat"), err)
	}

	stableStore, err := boltdb.NewBoltStore(filepath.Join(baseDir, "stable.dat"))
	if err != nil {
		return fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "stable.dat"), err)
	}

	snapshotStore, err := raft.NewFileSnapshotStore(baseDir, 3, os.Stderr)
	if err != nil {
		return fmt.Errorf(`raft.NewFileSnapshotStore(%q, ...): %v`, baseDir, err)
	}

	fsm := &fsm{store: s.store}

	maxPool := 5
	timeout := 10 * time.Second
	transport := raft.NewNetworkTransport(
		s.config.Raft.StreamLayer,
		maxPool,
		timeout,
		os.Stderr,
	)

	config := raft.DefaultConfig()
	config.LocalID = s.config.Raft.LocalID
	if s.config.Raft.HeartbeatTimeout != 0 {
		config.HeartbeatTimeout = s.config.Raft.HeartbeatTimeout
	}
	if s.config.Raft.ElectionTimeout != 0 {
		config.ElectionTimeout = s.config.Raft.ElectionTimeout
	}
	if s.config.Raft.LeaderLeaseTimeout != 0 {
		config.LeaderLeaseTimeout = s.config.Raft.LeaderLeaseTimeout
	}
	if s.config.Raft.CommitTimeout != 0 {
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
		return err
	}
	hasState, err := raft.HasExistingState(
		logStore,
		stableStore,
		snapshotStore,
	)
	if err != nil {
		return err
	}
	if s.config.Raft.Bootstrap && !hasState {
		config := raft.Configuration{
			Servers: []raft.Server{{
				ID:      config.LocalID,
				Address: raft.ServerAddress(s.config.Raft.BindAddr),
			}},
		}
		err = s.raft.BootstrapCluster(config).Error()
	}
	return err
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
	timeout := 10 * time.Second
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
	configFuture := s.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		return err
	}
	serverID := raft.ServerID(id)
	serverAddr := raft.ServerAddress(addr)
	for _, srv := range configFuture.Configuration().Servers {
		if srv.ID == serverID || srv.Address == serverAddr {
			if srv.ID == serverID && srv.Address == serverAddr {
				// server has already joined
				return nil
			}
			// remove the existing server
			removeFuture := s.raft.RemoveServer(serverID, 0, 0)
			if err := removeFuture.Error(); err != nil {
				return err
			}
		}
	}
	addFuture := s.raft.AddVoter(serverID, serverAddr, 0, 0)
	if err := addFuture.Error(); err != nil {
		return err
	}
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
		return err
	}
	return s.store.Close()
}

func (s *DistributedStorage) GetServers() ([]*api.Server, error) {
	future := s.raft.GetConfiguration()
	if err := future.Error(); err != nil {
		return nil, err
	}
	var servers []*api.Server
	for _, server := range future.Configuration().Servers {
		servers = append(servers, &api.Server{
			Id:       string(server.ID),
			RpcAddr:  string(server.Address),
			IsLeader: s.raft.Leader() == server.Address,
		})
	}
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
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := f.Write(body); err != nil {
		log.Fatal(err)
	}
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	fi, err := os.OpenFile(logFile, os.O_RDONLY, 0600)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer fi.Close()
	b, err := io.ReadAll(fi)
	if err != nil {
		log.Println(err)
		return nil, err
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
