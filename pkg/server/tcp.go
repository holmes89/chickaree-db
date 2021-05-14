package server

import (
	"fmt"
	"log"
	"net"

	"github.com/holmes89/chickaree-db/pkg/core"
)

type tcpServer struct {
	listener net.Listener
	repo     core.Repository
	errch    chan error
}

func NewTCPServer(port string, repo core.Repository) Runner {

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("listening on port %s\n", port)

	return &tcpServer{
		listener: listener,
		repo:     repo,
		errch:    make(chan error),
	}
}

func (s *tcpServer) Run() <-chan error {
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				s.errch <- err
			}
			_ = core.NewClient(conn, s.repo)
			fmt.Println("client connected")
		}
	}()
	return s.errch
}

func (s *tcpServer) Close() error {
	close(s.errch)
	return s.listener.Close()
}
