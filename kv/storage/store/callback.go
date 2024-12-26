package store

import (
	"AdachiAndShimamura/DistributedKV/proto/gen/raftpb"
	"errors"
)

type CallBack struct {
	done chan struct{}
	resp *raftpb.RaftCmdResponse
}

func NewCallBack() *CallBack {
	return &CallBack{
		done: make(chan struct{}, 1),
		resp: &raftpb.RaftCmdResponse{},
	}
}

func (cb *CallBack) Done() {
	cb.done <- struct{}{}
}

func (cb *CallBack) WaitResp() *raftpb.RaftCmdResponse {
	<-cb.done
	return cb.resp
}

func CheckCmdResp(resp *raftpb.RaftCmdResponse, reqCount int) error {
	if resp.Header == nil || len(resp.Resp) != reqCount {
		return errors.New("resp error")
	}
	return nil
}
