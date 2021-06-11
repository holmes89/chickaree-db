package storage

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

type event struct {
	Command string
	Args    [][]byte
}

type eventLogger struct {
	mu     sync.RWMutex
	events chan event
	pubs   map[chan event]bool
}

const logFile = "events.log"

func NewEventLogger() *eventLogger {
	evts := make(chan event)
	pubs := make(map[chan event]bool)
	return &eventLogger{
		mu:     sync.RWMutex{},
		events: evts,
		pubs:   pubs,
	}
}

func (l *eventLogger) Close() {
	close(l.events)
	for m := range l.pubs {
		close(m)
	}
}

func (l *eventLogger) Run() {
	for evt := range l.events {
		for m := range l.pubs {
			m <- evt
		}
	}
}

func (l *eventLogger) Leave(ch chan event) {
	delete(l.pubs, ch)
	close(ch)
}

func (l *eventLogger) Write(command string, args [][]byte) {
	l.mu.Lock()
	defer l.mu.Unlock()

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

	l.events <- event{
		Command: command,
		Args:    args,
	}
}

// May need to rethink this in the future, this will block writes as the file is read but if the file is large this will impact performance
func (l *eventLogger) Events() (chan event, error) {
	evts := make(chan event, 100)
	var buf []event

	defer l.mu.RLock()

	f, err := os.OpenFile(logFile, os.O_RDONLY, 0600)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	scanner := bufio.NewScanner(f)
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
	defer f.Close()

	l.pubs[evts] = true
	go l.eventLogs(evts, buf)

	return evts, nil
}

func (l *eventLogger) eventLogs(evts chan event, buf []event) {
	for _, evt := range buf {
		evts <- evt
	}
}
