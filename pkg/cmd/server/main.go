package main

import (
	"flag"
	"log"

	"github.com/holmes89/chickaree-db/pkg/core"
	"github.com/holmes89/chickaree-db/pkg/server"
)

func main() {

	var port string
	flag.StringVar(&port, "port", "6379", "port to listen on")

	flag.Parse()

	repo := core.NewRepo("chickaree.db")
	defer repo.Close()

	tcpServer := server.NewTCPServer(port, repo)
	defer tcpServer.Close()

	log.Println(<-tcpServer.Run())
	// grpcServer := server.NewGRPCServer(repo)
	// defer grpcServer.GracefulStop()

}
