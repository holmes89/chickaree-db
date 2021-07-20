package main

import (
	"flag"

	"github.com/holmes89/chickaree-db/chickaree"
	"github.com/holmes89/chickaree-db/chickaree/redis"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {

	var port string
	flag.StringVar(&port, "port", "6379", "port to listen on")
	flag.Parse()

	url := ":8080" // TODO discovery

	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(url, opts...)
	if err != nil {
		log.Fatal().Err(err).Str("url", url).Msg("failed to dial GRPC")
	}
	defer conn.Close()
	client := chickaree.NewChickareeDBClient(conn)

	tcpServer := redis.NewTCPServer(port, client)
	defer tcpServer.Close()

	log.Error().Err(<-tcpServer.Run()).Msg("terminated")
}
