package common

import "AdachiAndShimamura/DistributedKV/kv/engine_util"

type ClockConfig struct {
	HeartbeatInterval    uint64
	ElectionTimeoutStart uint64
	ElectionTimeoutEnd   uint64
	VoteTimeout          uint64
	BaseTimeInternal     uint64
}

type Config struct {
	//存储位置
	DBPath string

	Clock ClockConfig
}

type PeerConfig struct {
	Region   *Region
	ClockCfg *ClockConfig
	Storage  *engine_util.Engines
	Peers    map[uint64]string
}

type Region struct {
	RegionID uint64
	Start    []byte
	End      []byte
}
