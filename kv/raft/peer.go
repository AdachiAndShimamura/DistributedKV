package raft

import (
	. "AdachiAndShimamura/DistributedKV/kv/common"
	"AdachiAndShimamura/DistributedKV/kv/engine_util"
	"AdachiAndShimamura/DistributedKV/kv/rpc_client"
	"AdachiAndShimamura/DistributedKV/kv/utils"
	"AdachiAndShimamura/DistributedKV/proto/gen/raftpb"
)

type PeerMetaData struct {
	region   *Region
	clockCfg *ClockConfig
}

// Peer 表示Raft组的一个节点
type Peer struct {
	peerID   uint64
	regionID uint64
	storage  *engine_util.Engines
	meta     *PeerMetaData
	raft     *Raft
	log      *RaftLog
	client   *rpc_client.RaftClient
	msgs     chan *Msg
	tickerCh chan utils.TickType
	ticker   *utils.PeerTicker
}

func NewPeer(cfg *PeerConfig, peerID uint64, regionID uint64) *Peer {
	peer := &Peer{
		peerID:   peerID,
		regionID: regionID,
		meta: &PeerMetaData{
			region:   cfg.Region,
			clockCfg: cfg.ClockCfg,
		},
		tickerCh: make(chan utils.TickType, 2),
	}
	log := NewRaftLog(cfg.Storage, 0)
	raft := NewRaft(&RaftConfig{
		PeerID:  0,
		peers:   nil,
		applied: 0,
	}, log, peer)
	peer.log = log
	peer.raft = raft
	return peer
}
func (p *Peer) SendRaftMessage(peerID uint64,msg *raftpb.Message) error {
	return p.client.Send(peerID,&raftpb.RaftMessage{
		RegionId:    0,
		FromPeer:    &raftpb.Peer{
			Id:      peerID,
			StoreId: p.meta.region.RegionID,
		},
		ToPeer:      nil,
		Message:     nil,
		RegionEpoch: nil,
		StartKey:    nil,
		EndKey:      nil,
	})
}

func (p *Peer) Start() {
	p.raft.StartRaft()
	p.ticker.Register(utils.ElectionTimer, p.peerID, func() {
		p.tickerCh <- utils.ElectionTimer
	})
	for {
		select {
		case msg := <-p.msgs:
			{
				switch msg.Type {
				case MsgTypeRaftCmd:
					{
						p.raft.HandleRaftMessage((*raftpb.RaftMessage)(msg.Data))
					}
				case MsgTypeRaftMessage:
					{
						p.raft.
					}
				}
			}
		case msgType := <-p.tickerCh:
			{
				switch msgType {
				case utils.ElectionTimer:
				case utils.HeartBeat:

				}
			}
		}
	}
}
