package main

import (
	"flag"

	"github.com/holmes89/chickaree-db/bazel-chickaree-db/pkg/server"
	"github.com/holmes89/chickaree-db/pkg/core"
)

func main() {

	var port string
	flag.StringVar(&port, "port", "6379", "port to listen on")

	flag.Parse()

	repo := core.NewRepo("chickaree.db")
	defer repo.Close()

	tcpServer := server.NewTCPServer(port, repo)
	defer tcpServer.Close()

	grpcServer := server.NewGRPCServer(repo)
	defer grpcServer.GracefulStop()

}
