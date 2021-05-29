package redis

import (
	"fmt"
	"log"
	"net"
)

type TcpServer struct {
	listener net.Listener
	errch    chan error
}

func NewTCPServer(port string) *TcpServer {

	if port[0] != ':' {
		port = ":" + port
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("listening on port %s\n", port)

	errch := make(chan error)
	return &TcpServer{
		listener: listener,
		errch:    errch,
	}
}

func (s *TcpServer) Run() <-chan error {
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				s.errch <- err
			}
			_ = NewClient(conn)
			fmt.Println("client connected")
		}
	}()
	return s.errch
}

func (s *TcpServer) Close() error {
	log.Println("closing server...")
	close(s.errch)
	return s.listener.Close()
}
