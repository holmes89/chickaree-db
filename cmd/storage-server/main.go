package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/holmes89/chickaree-db/chickaree"
	"github.com/holmes89/chickaree-db/chickaree/storage"
	"google.golang.org/grpc"
)

func main() {
	port := 8080
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	gsrv := grpc.NewServer()

	srv, err := storage.NewServer(storage.ServerConfig{
		Config: storage.Config{
			StoragePath: "chickaree.db",
			RaftDir:     "/tmp",
		},
	})
	if err != nil {
		log.Panic(err)
	}
	defer srv.Close()
	chickaree.RegisterChickareeDBServer(gsrv, srv)

	errs := make(chan error, 2) // This is used to handle and log the reason why the application quit.
	go func() {
		fmt.Printf("listening on port %d...\n", port)
		errs <- gsrv.Serve(lis)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
		gsrv.GracefulStop()
	}()

	fmt.Printf("terminated: %s\n", <-errs)
}
