package server

import (
	"github.com/holmes89/chickaree-db/pkg/core"
	"github.com/holmes89/chickaree-db/pkg/proto"
	"google.golang.org/grpc"
)

type grpcServer struct {
	proto.UnimplementedChickareeDBServer
	repo core.Repository
}

func NewGRPCServer(repo core.Repository) *grpc.Server {
	gsrv := grpc.NewServer()
	srv := &grpcServer{
		repo: repo,
	}
	proto.RegisterChickareeDBServer(gsrv, srv)
	return gsrv
}
