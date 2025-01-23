package rpc_server

import (
	"AdachiAndShimamura/DistributedKV/kv/storage"
	rpcpb "AdachiAndShimamura/DistributedKV/proto/gen"
	"log/slog"
)

type Server struct {
	storage *storage.RaftStorage
	rpcpb.UnimplementedTinyKvRpcServer
}

func NewServer(config *storage.Config) *Server {
	s := storage.NewRaftStorage(config)
	return &Server{
		storage: s,
	}
}

func (s *Server) Raft(stream rpcpb.TinyKvRpc_RaftServer) error {
	err := s.storage.HandleRaftMessage(&stream)
	if err != nil {
		slog.Info("err:", err)
	}
	return nil
}
