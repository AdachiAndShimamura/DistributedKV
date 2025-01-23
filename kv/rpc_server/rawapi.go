package rpc_server

import (
	"AdachiAndShimamura/DistributedKV/kv/common"
	"AdachiAndShimamura/DistributedKV/proto/gen/kvrpcpb"
	"AdachiAndShimamura/DistributedKV/proto/gen/raftpb"
	"context"
)

func (s *Server) RawGet(ctx context.Context, in *kvrpcpb.RawGet) (*kvrpcpb.RawGetResponse, error) {
	cmd := &raftpb.RaftCmdRequest{
		Header: &raftpb.RequestHeader{
			RegionId: in.Ctx.RegionId,
		},
		Requests: []*raftpb.Request{{
			CmdType: raftpb.RaftCmdType_Get,
			Get: &raftpb.GetReq{
				Key: in.Key,
				Cf:  in.Cf,
			},
		}},
	}
	cb := common.NewCallBack()
	err := s.storage.HandleCmdMessage(cmd, cb)
	if err != nil {
		return nil, err
	}
	resp := cb.WaitResp()
	if err = common.CheckCmdResp(resp, 1); err != nil {
		return nil, err
	}
	return &kvrpcpb.RawGetResponse{
		Value: resp.Resp[0].Get.GetVal(),
		Have:  true,
		Error: "",
	}, nil
}

func (s *Server) RawPut(ctx context.Context, in *kvrpcpb.RawPut) (*kvrpcpb.RawPutResponse, error) {
	return nil, nil
}

func (s *Server) RawDelete(ctx context.Context, in *kvrpcpb.RawDelete) (*kvrpcpb.RawDeleteResponse, error) {
	return nil, nil
}

func (s *Server) RawScan(ctx context.Context, in *kvrpcpb.RawScan) (*kvrpcpb.RawScanResponse, error) {
	return nil, nil
}
