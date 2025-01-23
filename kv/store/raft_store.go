package store

import (
	. "AdachiAndShimamura/DistributedKV/kv/common"
	"AdachiAndShimamura/DistributedKV/kv/raft"
	"AdachiAndShimamura/DistributedKV/kv/utils"
)

type RaftStoreMetaData struct {
	clockCfg *ClockConfig
}

// RaftStore 表示作为存储节点的一台机器，以StoreID标识每台机器
type RaftStore struct {
	storeID  uint64
	metaData *RaftStoreMetaData
	router   *Router
	regions  map[uint64]*raft.Peer
	ticker   *utils.TickerCore
}

func CreateRaftStore(config *Config, storeID uint64, router *Router) *RaftStore {
	md := &RaftStoreMetaData{
		clockCfg: &ClockConfig{
			HeartbeatInterval:    config.HeartbeatInterval,
			ElectionTimeoutStart: config.ElectionTimeoutStart,
			ElectionTimeoutEnd:   config.ElectionTimeoutEnd,
			BaseTimeInternal:     config.BaseTimeInternal,
		},
	}
	return &RaftStore{storeID: storeID, router: router, metaData: md}
}

func (r *RaftStore) Start() error {
	utils.TickerManager.Start(r.metaData.clockCfg)
	for _, peer := range r.regions {
		go peer.Start()
	}
	return nil
}

func (r *RaftStore) LoadPeers() error {
	cfg := &PeerConfig{
		Region: &Region{
			RegionID: 1,
			Start:    []byte{},
			End:      nil,
		},
		ClockCfg: r.metaData.clockCfg,
	}
	peer1 := raft.NewPeer(cfg, 1, 1)
	r.regions[1] = peer1
	return nil
}

func (r *RaftStore) AddTimer() {

}

func (r *RaftStore) Stop() {
}
