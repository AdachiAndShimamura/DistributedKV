package kv

import (
	"AdachiAndShimamura/DistributedKV/proto/gen/kvrpcpb"
	"AdachiAndShimamura/DistributedKV/proto/gen/raftpb"
	"context"
)

func (s *Server) RawGet(ctx context.Context, in *kvrpcpb.RawGet) (*kvrpcpb.RawGetResponse, error) {
	cmd := &raftpb.RaftCmdRequest{
		Ctx:     in.Ctx,
		CmdType: raftpb.RaftCmdType_Get,
		Key:     in.Key,
		Value:   nil,
	}
	s.storage.HandleCmdMessage(cmd)
}

func (s *Server) RawPut(ctx context.Context, in *kvrpcpb.RawPut) (*kvrpcpb.RawPutResponse, error) {

}

func (s *Server) RawDelete(ctx context.Context, in *kvrpcpb.RawDelete) (*kvrpcpb.RawDeleteResponse, error) {

}

func (s *Server) RawScan(ctx context.Context, in *kvrpcpb.RawScan) (*kvrpcpb.RawScanResponse, error) {

}
