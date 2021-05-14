package server

import (
	"github.com/holmes89/chickaree-db/pkg/core"
	"github.com/holmes89/chickaree-db/pkg/proto"
	"google.golang.org/grpc"
)

type grpcServer struct {
	proto.UnimplementedReplicatorServer
	repo core.Repository
}

func NewGRPCServer(repo core.Repository) *grpc.Server {
	gsrv := grpc.NewServer()
	srv := &grpcServer{
		repo: repo,
	}
	proto.RegisterReplicatorServer(gsrv, srv)
	return gsrv
}

func (s *grpcServer) Subscribe(sub proto.Replicator_SubscribeServer) error {
	for {
		req, err := sub.Recv()
		if err != nil {
			return err
		}
		r := NewRequest(req.Request)
		if r == nil {
			continue
		}
		s.repo.Handle(*r)
	}
}

func NewRequest(req *proto.Request) *core.Request {
	if req == nil {
		return nil
	}
	count := int(req.MsgCount)
	var args []core.Arg
	for _, r := range req.Args {
		args = append(args, core.Arg(r))
	}

	resp := core.Request{
		MsgCount: count,
		Args:     args,
		Command:  req.Command,
	}
	return &resp
}
