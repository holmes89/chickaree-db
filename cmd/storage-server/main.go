package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/holmes89/chickaree-db/chickaree"
	"github.com/holmes89/chickaree-db/chickaree/storage"
	"github.com/rs/zerolog/log"
	"github.com/soheilhy/cmux"
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

	mux := srv.Mux()
	grpcLn := mux.Match(cmux.Any())
	chickaree.RegisterChickareeDBServer(gsrv, srv)

	errs := make(chan error, 2) // This is used to handle and log the reason why the application quit.
	go func() {
		log.Info().Int("port", cfg.RPCPort).Msg("listening grpc...")
		errs <- gsrv.Serve(grpcLn)
	}()
	go func() {
		log.Info().Int("port", cfg.RPCPort).Msg("listening mux...")
		errs <- mux.Serve()
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
		gsrv.GracefulStop()
	}()

	log.Error().Err(<-errs).Msg("terminated")
}
