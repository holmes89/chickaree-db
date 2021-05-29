package main

import (
	"flag"
	"log"

	"github.com/holmes89/chickaree-db/pkg/core/redis"
)

func main() {

	var port string
	flag.StringVar(&port, "port", "6379", "port to listen on")

	flag.Parse()

	tcpServer := redis.NewTCPServer(port)
	defer tcpServer.Close()

	log.Println(<-tcpServer.Run())
	// grpcServer := server.NewGRPCServer(repo)
	// defer grpcServer.GracefulStop()

}
