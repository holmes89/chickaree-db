package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/holmes89/chickaree-db/chickaree"
	"github.com/holmes89/chickaree-db/chickaree/storage"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {

	gsrv := grpc.NewServer()

	cfg, err := storage.LoadConfiguration()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to load configuration")
	}
	srv, err := storage.NewServer(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create server")
	}
	defer srv.Close()
	chickaree.RegisterChickareeDBServer(gsrv, srv)

	errs := make(chan error, 2) // This is used to handle and log the reason why the application quit.
	go func() {
		log.Info().Int("port", cfg.RPCPort).Msg("listening...")
		errs <- gsrv.Serve(srv.Mux())
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
		gsrv.GracefulStop()
	}()

	log.Error().Err(<-errs).Msg("terminated")
}
