package server

import (
	"fmt"
	"log"
	"net"

	"github.com/holmes89/chickaree-db/pkg/core"
	"github.com/holmes89/chickaree-db/pkg/core/redis"
)

type tcpServer struct {
	listener net.Listener
	repo     core.Repository
	errch    chan error
}

func NewTCPServer(port string, repo core.Repository) Runner {

	if port[0] != ':' {
		port = ":" + port
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("listening on port %s\n", port)

	errch := make(chan error)
	return &tcpServer{
		listener: listener,
		repo:     repo,
		errch:    errch,
	}
}

func (s *tcpServer) Run() <-chan error {
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				s.errch <- err
			}
			_ = redis.NewClient(conn, s.repo)
			fmt.Println("client connected")
		}
	}()
	return s.errch
}

func (s *tcpServer) Close() error {
	log.Println("closing server...")
	close(s.errch)
	return s.listener.Close()
}
