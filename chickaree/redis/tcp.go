package redis

import (
	"context"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/holmes89/chickaree-db/chickaree"
)

type TcpServer struct {
	listener     net.Listener
	client       chickaree.ChickareeDBClient
	leaderConn   *grpc.ClientConn
	leaderClient chickaree.ChickareeDBClient
	leaderURL    string
	errch        chan error
	done         chan bool
	ticker       *time.Ticker
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
	s := &TcpServer{
		listener: listener,
		errch:    errch,
		client:   client,
	}
	s.setLeaderClient()
	s.ticker = time.NewTicker(30 * time.Second)
	s.done = make(chan bool)
	go func() {
		for {
			select {
			case <-s.done:
				return
			case <-s.ticker.C:
				s.setLeaderClient()
			}
		}
	}()

	return s
}

func (s *TcpServer) Run() <-chan error {
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				s.errch <- err
			}
			_ = NewClient(conn, s.leaderClient, s.client)
			log.Info().Msg("client connected")
		}
	}()
	return s.errch
}

func (s *TcpServer) Close() error {
	log.Info().Msg("closing server...")
	close(s.errch)
	s.ticker.Stop()
	s.done <- true
	close(s.done)
	s.leaderConn.Close()
	return s.listener.Close()
}

func (s *TcpServer) setLeaderClient() {
	log.Info().Msg("setting leader client...")
	resp, err := s.client.GetServers(context.Background(), &chickaree.GetServersRequest{})
	if err != nil || resp == nil || resp.Servers == nil {
		log.Error().Err(err).Msg("unable to find servers")
		return
	}
	for _, sv := range resp.Servers {
		if sv.IsLeader {
			if sv.RpcAddr == s.leaderURL {
				log.Info().Str("url", sv.RpcAddr).Msg("leader has not changed.")
				return
			}
			opts := []grpc.DialOption{grpc.WithInsecure()}
			conn, err := grpc.Dial(sv.RpcAddr, opts...)
			log.Info().Str("url", sv.RpcAddr).Msg("dialing leader...")
			if err != nil {
				log.Fatal().Err(err).Str("url", sv.RpcAddr).Msg("failed to dial leader GRPC")
				return
			}
			if s.leaderConn != nil {
				if err := s.leaderConn.Close(); err != nil {
					log.Error().Err(err).Msg("error closing leader conn")
				}
			}

			s.leaderConn = conn
			s.leaderClient = chickaree.NewChickareeDBClient(conn)
			log.Info().Msg("leader client set.")
			return
		}
	}
	log.Error().Msg("unable to find leader")
}
