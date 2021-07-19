package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/hashicorp/raft"
)

var (
	_ raft.FSM = &eventLogger{}
)

type event struct {
	Command string
	Args    [][]byte
}

type eventLogger struct {
	mu    sync.RWMutex
	store storage
}

const logFile = "events.log"

func NewEventLogger(store storage) *eventLogger {
	return &eventLogger{
		mu:    sync.RWMutex{},
		store: store,
	}
}

func (l *eventLogger) Close() error {
	return l.store.Close()
}

func (s *eventLogger) handleEvent(evt event) error {
	command := evt.Command
	args := evt.Args

	switch strings.ToLower(command) {
	case "set":
		return s.store.Set(args[0], args[1])
	default:
		return fmt.Errorf("unknown command: '%s'", command)
	}
}

func (l *eventLogger) Apply(lg *raft.Log) interface{} {
	var evt event
	if err := json.Unmarshal(lg.Data, &evt); err != nil {
		return err
	}
	if err := l.handleEvent(evt); err != nil {
		return err
	}

	return evt
}

func (l *eventLogger) Snapshot() (raft.FSMSnapshot, error) {
	// Make sure that any future calls to f.Apply() don't change the snapshot.
	return l, nil
}

func (l *eventLogger) Restore(r io.ReadCloser) error {
	evts := readEvents(r)
	for _, evt := range evts {
		if err := l.handleEvent(evt); err != nil {
			return err
		}
	}
	return nil
}

func (l *eventLogger) Persist(sink raft.SnapshotSink) error {
	l.mu.RLock()
	defer l.mu.RUnlock()
	f, err := os.OpenFile(logFile, os.O_RDONLY, 0600)
	if err != nil {
		log.Println(err)
		return err
	}

	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = sink.Write(b)
	if err != nil {
		sink.Cancel()
		return fmt.Errorf("sink.Write(): %v", err)
	}
	return sink.Close()
}

func (l *eventLogger) Release() {
}

func (l *eventLogger) write(evt event) {
	l.mu.Lock()
	defer l.mu.Unlock()

	command := evt.Command
	args := evt.Args

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	eLog := fmt.Sprintf("%s", command)
	for _, arg := range args {
		eLog += fmt.Sprintf("\t%s", string(arg))
	}
	eLog += "\n"
	if _, err = f.WriteString(eLog); err != nil {
		log.Fatal(err)
	}
}

func readEvents(r io.ReadCloser) []event {
	var buf []event

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		evt := event{}
		for _, v := range strings.Split("\t", line) {
			if evt.Command == "" {
				evt.Command = v
				continue
			}
			evt.Args = append(evt.Args, []byte(v))
		}
		buf = append(buf, evt)
	}

	return buf
}
