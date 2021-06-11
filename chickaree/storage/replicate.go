package storage

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/holmes89/chickaree-db/chickaree"
	"google.golang.org/grpc"
)

type Replicator struct {
	DialOptions []grpc.DialOption
	Local       chickaree.ChickareeDBClient

	mu      sync.Mutex
	servers map[string]chan struct{}
	closed  bool
	close   chan struct{}
}

func (s *Replicator) init() {
	if s.servers == nil {
		s.servers = make(map[string]chan struct{})
	}
	if s.close == nil {
		s.close = make(chan struct{})
	}
}

func (s *Replicator) Join(name, addr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.init()

	if s.closed {
		return nil
	}

	if _, ok := s.servers[name]; ok {
		return nil // already replicating
	}
	s.servers[name] = make(chan struct{})

	go s.getEvents(addr, s.servers[name])
	return nil
}

func (s *Replicator) getEvents(addr string, leave chan struct{}) {
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("unable to get events from %s", addr)
	}
	defer cc.Close()

	client := chickaree.NewChickareeDBClient(cc)

	ctx := context.Background()
	stream, err := client.EventLog(ctx, &chickaree.EventLogRequest{})
	if err != nil {
		log.Fatalf("unable to get event stream from %s", addr)
	}

	events := make(chan event)
	defer close(events)
	go func() {
		for {
			recv, err := stream.Recv()
			if err != nil {
				log.Fatalf("unable to get event from stream from %s", addr)
			}
			events <- event{
				Command: string(recv.Command),
				Args:    recv.Args,
			}
		}
	}()

	for {
		select {
		case <-s.close:
			return
		case <-leave:
			return
		case evt := <-events:
			s.handleEvent(evt.Command, evt.Args)
		}
	}
}

func (s *Replicator) handleEvent(command string, args [][]byte) {
	switch strings.ToLower(command) {
	case "set":
		s.Local.Set(context.Background(), &chickaree.SetRequest{
			Key:   string(args[0]),
			Value: args[1],
		})
	default:
		log.Printf("unknown command '%s'", command)
	}
}

func (s *Replicator) Leave(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.init()

	if _, ok := s.servers[name]; !ok {
		return nil
	}

	close(s.servers[name])
	delete(s.servers, name)
	return nil
}

func (s *Replicator) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.init()

	if s.closed {
		return nil
	}

	s.closed = true
	close(s.close)
	return nil
}
