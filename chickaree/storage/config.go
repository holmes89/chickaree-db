package storage

import "github.com/hashicorp/raft"

type Config struct {
	StoragePath string
	RaftDir     string
	Raft        struct {
		raft.Config
		BindAddr    string
		StreamLayer *StreamLayer
		Bootstrap   bool
	}
}
