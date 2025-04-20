package raft

import (
	"AdachiAndShimamura/DistributedKV/kv/engine_util"
	"AdachiAndShimamura/DistributedKV/proto/gen/raftpb"
)

type RaftLog struct {
	// storage contains all stable entries since the last snapshot.
	storage *engine_util.Engines

	index uint64

	// committed is the highest log position that is known to be in
	// stable storage on a quorum of nodes.
	committed uint64

	// applied is the highest log position that the application has
	// been instructed to apply to its state machine.
	// Invariant: applied <= committed
	applied uint64

	// log entries with index <= stabled are persisted to storage.
	// It is used to record the logs that are not persisted by storage yet.
	// Everytime handling `Ready`, the unstabled logs will be included.
	stabled uint64

	currentTerm uint64

	// all entries that have not yet compact.
	entries []raftpb.Entry

	//// the incoming unstable snapshot, if any.
	//// (Used in 2C)
	//pendingSnapshot *pb.Snapshot

	// Your Data Here (2A).
}

func NewRaftLog(storage *engine_util.Engines, applied uint64) *RaftLog {
	return &RaftLog{
		storage:   storage,
		committed: applied,
		applied:   applied,
		stabled:   applied,
		entries:   make([]raftpb.Entry, 10),
	}
}
