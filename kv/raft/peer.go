package raft

import (
	"AdachiAndShimamura/DistributedKV/kv/storage/store"
	"AdachiAndShimamura/DistributedKV/kv/utils"
)

// Peer 表示Raft组的一个节点
type Peer struct {
	raft   *Raft
	log    *RaftLog
	msgs   chan store.Msg
	ticker *utils.PeerTicker
}

func NewPeer(raft *Raft, log *RaftLog) *Peer {

}

func (p *Peer) Start() {

}
