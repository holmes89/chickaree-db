package redis

import (
	"net"

	"github.com/rs/zerolog/log"

	"github.com/holmes89/chickaree-db/chickaree"
)

type TcpServer struct {
	listener net.Listener
	client   chickaree.ChickareeDBClient
	errch    chan error
}

func NewTCPServer(port string, client chickaree.ChickareeDBClient) *TcpServer {

	if port[0] != ':' {
		port = ":" + port
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal().Err(err).Str("port", port).Msg("failed to listen")
	}

	log.Info().Str("port", port).Msg("listening...")

	errch := make(chan error)
	return &TcpServer{
		listener: listener,
		errch:    errch,
		client:   client,
	}
}

func (s *TcpServer) Run() <-chan error {
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				s.errch <- err
			}
			_ = NewClient(conn, s.client)
			log.Info().Msg("client connected")
		}
	}()
	return s.errch
}

func (s *TcpServer) Close() error {
	log.Info().Msg("closing server...")
	close(s.errch)
	return s.listener.Close()
}
