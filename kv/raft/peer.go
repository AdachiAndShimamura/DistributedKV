package raft

import (
	. "AdachiAndShimamura/DistributedKV/kv/common"
	"AdachiAndShimamura/DistributedKV/kv/engine_util"
	"AdachiAndShimamura/DistributedKV/kv/rpc_client"
	"AdachiAndShimamura/DistributedKV/kv/utils"
	"AdachiAndShimamura/DistributedKV/proto/gen/raftpb"
	"log/slog"
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
	msgCh    chan *Msg
	tickerCh chan utils.ClockType
	ticker   *utils.PeerTicker
}

func NewPeer(cfg *PeerConfig, peerID uint64, regionID uint64) *Peer {
	ch := make(chan utils.ClockType)
	ticker := utils.NewPeerTicker(peerID, cfg.ClockCfg, ch)
	peer := &Peer{
		peerID:   peerID,
		regionID: regionID,
		meta: &PeerMetaData{
			region:   cfg.Region,
			clockCfg: cfg.ClockCfg,
		},
		ticker:   ticker,
		tickerCh: ch,
	}
	log := NewRaftLog(cfg.Storage, 0)
	raft := NewRaft(&RaftConfig{
		peerID:   peerID,
		regionID: regionID,
		peers:    nil,
		applied:  0,
	}, log, peer, ticker)
	peer.log = log
	peer.raft = raft
	return peer
}
func (p *Peer) SendRaftMessage(toPeer uint64, msg *raftpb.Message) {
	err := p.client.Send(toPeer, &raftpb.RaftMessage{
		RegionId: 0,
		FromPeer: &raftpb.Peer{
			Id:      p.peerID,
			StoreId: p.meta.region.RegionID,
		},
		ToPeer: &raftpb.Peer{
			Id: toPeer, StoreId: p.meta.region.RegionID,
		},
		Message:     msg,
		RegionEpoch: nil,
		StartKey:    nil,
		EndKey:      nil,
	})
	if err != nil {
		slog.Warn("send raft message to peer failed, peerID: %d, err: %v", p.peerID, err)
	}
}

func (p *Peer) Start() {
	p.raft.StartRaft()
	for {
		select {
		case msg := <-p.msgCh:
			{
				switch msg.Type {
				case MsgTypeRaftCmd:
					{
						raftMsg, ok := msg.Data.(raftpb.RaftMessage)
						if !ok {
							slog.Warn("raft msg err")
						}
						go p.raft.HandleRaftMessage(&raftMsg)
					}
				case MsgTypeRaftMessage:
					{
						cmdMsg, ok := msg.Data.(RaftCmdMsg)
						if !ok {
							slog.Warn("raft cmd msg err")
						}
						go p.raft.HandleCmdMessage(&cmdMsg)
					}
				}
			}
		case msgType := <-p.tickerCh:
			{
				switch msgType {
				case utils.ElectionTimerClock:
					go p.raft.HandleRaftMessage(&raftpb.RaftMessage{
						Message: &raftpb.Message{
							MsgType: raftpb.MessageType_MsgHup,
						},
					})
				case utils.HeartBeatClock:
					go p.raft.HandleRaftMessage(&raftpb.RaftMessage{
						Message: &raftpb.Message{
							MsgType: raftpb.MessageType_MsgBeat,
						},
					})
				}
			}
		}
	}
}

func (p *Peer) ResetClock(tp utils.ClockType) {
	p.ticker.Reset(tp)
}

func (p *Peer) Stop() {
	p.raft.peer = nil
}
