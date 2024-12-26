package raft

import (
	"AdachiAndShimamura/DistributedKV/kv/storage"
	rpc "AdachiAndShimamura/DistributedKV/proto/gen/raftpb"
	"context"
	"math/rand/v2"
	"sync"
	"time"
)

type NodeState int

const (
	Follower NodeState = iota
	Candidate
	Leader
)

type Config struct {
	ID      uint64
	peers   []uint64
	storage *storage.Storage
	applied uint64
	//heartbeatTick int
	//electionTick  int
}

// LogEntry 表示Raft协议中的一个日志条目
type LogEntry struct {
	Term    int
	Command interface{}
}
type Progress struct {
	match uint64
	next  uint64
}

// Raft 结构体，表示Raft节点
type Raft struct {
	ID          uint64
	State       NodeState
	currentTerm uint64
	votedFor    uint64
	LeaderId    uint64
	progress    map[uint64]*Progress
	mu          sync.Mutex
	log         *RaftLog

	//random,150ms-300ms
	electionTimeout int
	electionTicker  *time.Ticker

	electionTaskTimeout int
	electionTaskTimer   *time.Ticker

	//only leader
	heartbeatSend       int
	heartbeatSendTicker *time.Ticker
}

// NewRaft 创建一个Raft实例
func NewRaft(config Config) *Raft {
	progress := make(map[uint64]*Progress)
	for _, peer := range config.peers {
		progress[peer] = new(Progress)
	}
	electionTimeout := rand.IntN(150) + 150
	return &Raft{
		ID:              config.ID,
		State:           Follower,
		currentTerm:     1,
		votedFor:        0,
		progress:        progress,
		log:             NewRaftLog(config.storage, config.applied),
		electionTimeout: electionTimeout,
		heartbeatSend:   electionTimeout / 2,
	}
}

func (r *Raft) HandleRaftMessage(data *rpc.RaftMessage) {
	msg := data.Message
	msgType := msg.MsgType
	switch msgType {
	//本地消息：选举时钟结束，立即选举
	case rpc.MessageType_MsgHup:
		{

		}
	//本地消息：心跳时钟结束，立即发送心跳
	case rpc.MessageType_MsgBeat:
		{

		}
	case rpc.MessageType_MsgPropose:
		{

		}
	case rpc.MessageType_MsgAppend:
		{

		}
	case rpc.MessageType_MsgAppendResponse:
		{

		}
	case rpc.MessageType_MsgRequestVote:
		{

		}
	case rpc.MessageType_MsgRequestVoteResponse:
		{

		}
	case rpc.MessageType_MsgSnapshot:
		{

		}
	case rpc.MessageType_MsgHeartbeat:
		{

		}
	case rpc.MessageType_MsgHeartbeatResponse:
		{

		}
	case rpc.MessageType_MsgTransferLeader:
		{

		}
	case rpc.MessageType_MsgTimeoutNow:
		{

		}
	}
}

func (r *Raft) StartRaft() {
	r.electionTicker = time.NewTicker(time.Duration(r.electionTimeout) * time.Millisecond)
	r.heartbeatTicker = time.NewTicker(time.Duration(r.heartbeatTimeout) * time.Millisecond)
	ctx := context.Background()
	for {
		select {
		case <-r.electionTicker.C:
			r.StartElection(ctx)
		case <-r.heartbeatTicker.C:

		}
	}
}

// StartElection 发起选举请求
func (r *Raft) StartElection(ctx context.Context) {

}

// AppendEntries 处理日志条目追加
func (r *Raft) AppendEntries(term int, entries []LogEntry) {

}

// ApplyLogs 应用已提交的日志
func (r *Raft) ApplyLogs() {

}

// Heartbeat 领导者发送心跳信号，保持自己的领导者身份
func (r *Raft) Heartbeat() {
	r.mu.Lock()
	defer r.mu.Unlock()
	// 如果是领导者，定期发送心跳
	if r.state == Leader {
		// 发送心跳给所有跟随者
		// 在真实实现中这里需要 RPC
	}
}
func (r *Raft) AddNode(id uint64) {

}

// StateMachine 用于处理具体的业务逻辑（例如KV存储的操作）
type StateMachine struct {
	data map[string]interface{}
}

// NewStateMachine 创建一个状态机实例
func NewStateMachine() *StateMachine {
	return &StateMachine{
		data: make(map[string]interface{}),
	}
}

// ApplyCommand 应用命令到状态机（KV存储的执行）
func (sm *StateMachine) ApplyCommand(command interface{}) {
	switch cmd := command.(type) {
	case PutCommand:
		sm.data[cmd.Key] = cmd.Value
	case GetCommand:
		// 返回存储的值
	}
}

// PutCommand 和 GetCommand 用于模拟KV操作
type PutCommand struct {
	Key   string
	Value interface{}
}

type GetCommand struct {
	Key string
}
