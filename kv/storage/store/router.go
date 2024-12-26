package store

import (
	"AdachiAndShimamura/DistributedKV/proto/gen/raftpb"
	"github.com/pkg/errors"
	"sync"
)

// Router 将Message路由到对应的Peer，交由对应的raft状态机处理
type Router struct {
	//grpc是多线程服务，因此peers映射也应该支持多线程下的读写访问
	peerSender  sync.Map
	storeSender chan *Msg
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) send(regionId uint64, msg *Msg) error {
	data, ok := r.peerSender.Load(regionId)
	if !ok {
		return errors.New("peer not exist")
	}
	ch := data.(chan *Msg)
	ch <- msg
	return nil
}

func (r *Router) sendStore(msg *Msg) {
	r.storeSender <- msg
}

func (r *Router) SendRaftMessage(msg *raftpb.RaftMessage) {
	id := msg.RegionId
	if err := r.send(id, NewStoreMsg(MsgTypeRaftMessage, id, msg)); err != nil {
		r.sendStore(NewStoreMsg(MsgTypeStoreRaftMessage, id, msg))
	}
}

func (r *Router) SendRaftCmdMessage(msg *raftpb.RaftCmdRequest, cb *CallBack) error {
	cmd := &RaftCmdMsg{
		cb:  cb,
		req: msg,
	}
	return r.send(msg.Header.RegionId, NewStoreMsg(MsgTypeRaftCmd, msg.Header.RegionId, cmd))
}
