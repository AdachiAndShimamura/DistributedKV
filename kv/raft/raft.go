package raft

import (
	. "AdachiAndShimamura/DistributedKV/kv/common"
	"AdachiAndShimamura/DistributedKV/kv/utils"
	"AdachiAndShimamura/DistributedKV/proto/gen/raftpb"
	"context"
	"log/slog"
	"sync"
	"time"
)

type NodeState int

const (
	Follower NodeState = iota
	Candidate
	Leader
)

type RaftConfig struct {
	peerID     uint64
	peers      []uint64
	regionID   uint64
	applied    uint64
	selfPeer   *Peer
	timeConfig *ClockConfig
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
type name struct {
}

// Raft 结构体，表示Raft节点
type Raft struct {
	PeerID      uint64
	RegionId    uint64
	State       NodeState
	currentTerm uint64
	votedFor    uint64
	LeaderId    uint64
	progress    map[uint64]*Progress
	mu          sync.Mutex
	log         *RaftLog
	peer        *Peer
	peers       []uint64
	es          electionStatus
	timeCfg     *ClockConfig
	timer       *utils.PeerTicker
}

// NewRaft 创建一个Raft实例
func NewRaft(config *RaftConfig, log *RaftLog, peer *Peer, ticker *utils.PeerTicker) *Raft {
	progress := make(map[uint64]*Progress)
	for _, peer := range config.peers {
		progress[peer] = new(Progress)
	}
	return &Raft{
		PeerID:      config.peerID,
		RegionId:    config.regionID,
		State:       Follower,
		currentTerm: 1,
		votedFor:    0,
		progress:    progress,
		log:         log,
		peer:        peer,
		peers:       config.peers,
		timer:       ticker,
	}
}

func (r *Raft) HandleRaftMessage(data *raftpb.RaftMessage) {
	msg := data.Message
	msgType := msg.MsgType
	switch r.State {
	case Follower:
		{
			switch msgType {
			//本地消息：选举时钟结束，立即选举
			case raftpb.MessageType_MsgHup:
				{
					go r.StartElection()
				}
			case raftpb.MessageType_MsgAppend:
				{
					//若任期小于自己，拒绝，返回自己的任期
					if msg.Term < r.currentTerm {
						go r.peer.SendRaftMessage(msg.From, &raftpb.Message{
							MsgType: 0,
							To:      0,
							From:    0,
							Term:    0,
							LogTerm: 0,
							Index:   0,
							Entries: nil,
							Commit:  0,
							Reject:  false,
						})
					}
					//若日志与自身不一致，删除多余的日志
					if
				}
			case raftpb.MessageType_MsgRequestVote:
				{
					r.mu.Lock()
					//这里没有考虑有没有可能出现一个Term很高，但是在任期内没有发起过投票的Follower
					if r.votedFor == msg.From || (r.votedFor == 0 && msg.Term > r.currentTerm && msg.LogTerm >= r.log.currentTerm && msg.Index >= r.log.index) {
						//承诺接受
						r.votedFor = msg.From
						r.currentTerm = msg.Term
						go r.peer.SendRaftMessage(msg.From, &raftpb.Message{
							MsgType: raftpb.MessageType_MsgRequestVoteResponse,
							To:      msg.From,
							From:    r.PeerID,
							Term:    r.currentTerm,
							LogTerm: r.log.currentTerm,
							Index:   r.log.index,
							Entries: nil,
							Commit:  r.log.committed,
							Reject:  false,
						})
					} else {
						go r.peer.SendRaftMessage(msg.From, &raftpb.Message{
							MsgType: raftpb.MessageType_MsgRequestVoteResponse,
							To:      msg.From,
							From:    r.PeerID,
							Term:    r.currentTerm,
							LogTerm: r.log.currentTerm,
							Index:   r.log.index,
							Entries: nil,
							Commit:  r.log.committed,
							Reject:  true,
						})
					}
					r.mu.Unlock()
				}
			case raftpb.MessageType_MsgSnapshot:
				{
				}
			case raftpb.MessageType_MsgHeartbeat:
				{

				}
			case raftpb.MessageType_MsgTimeoutNow:
				{
				}
			default:
				{
					slog.Info("Unknown raft message")
				}
			}
		}
	case Candidate:
		{
			switch msgType {
			case raftpb.MessageType_MsgBeat:
				{

				}
			case raftpb.MessageType_MsgAppend:
				{

				}
			case raftpb.MessageType_MsgAppendResponse:
				{

				}
			case raftpb.MessageType_MsgRequestVote:
				{
					r.mu.Lock()
					defer r.mu.Unlock()
					//Candidate的Term要大于自身，并且日志完整，才同意投票
					if msg.Term > r.currentTerm && msg.Index >= r.log.index {
						r.currentTerm = msg.Term
						r.votedFor = msg.From
						r.becomeFollower()
						go r.peer.SendRaftMessage(msg.From, &raftpb.Message{
							MsgType: raftpb.MessageType_MsgRequestVoteResponse,
							To:      msg.From,
							From:    r.PeerID,
							Term:    r.currentTerm,
							LogTerm: r.log.currentTerm,
							Index:   r.log.index,
							Entries: nil,
							Commit:  r.log.committed,
							Reject:  false,
						})
					} else {
						go r.peer.SendRaftMessage(msg.From, &raftpb.Message{
							MsgType: raftpb.MessageType_MsgRequestVoteResponse,
							To:      msg.From,
							From:    r.PeerID,
							Term:    r.currentTerm,
							LogTerm: r.log.currentTerm,
							Index:   r.log.index,
							Entries: nil,
							Commit:  r.log.committed,
							Reject:  true,
						})
					}
				}
			case raftpb.MessageType_MsgRequestVoteResponse:
				{
					if !msg.Reject {
						r.es.statistics(msg.From, msg.Reject)
					} else {
						//Candidate后处理
						//若是遇到更大的Term,立刻退回到Follower（这里要防止掉线节点无限增加Term,要引入PreVote）
						if msg.Term > r.currentTerm {
							r.currentTerm = msg.Term
							r.votedFor = 0
							r.becomeFollower()
						} else if msg.Term == r.currentTerm && msg.From == r.LeaderId {
							r.becomeFollower()
						}
					}

				}
			case raftpb.MessageType_MsgSnapshot:
				{

				}
			case raftpb.MessageType_MsgHeartbeat:
				{

				}
			case raftpb.MessageType_MsgHeartbeatResponse:
				{

				}
			case raftpb.MessageType_MsgTransferLeader:
				{

				}
			case raftpb.MessageType_MsgTimeoutNow:
				{

				}
			}
		}
	case Leader:
		{
			switch msgType {
			//本地消息：心跳时钟结束，立即发送心跳
			case raftpb.MessageType_MsgBeat:
				{

				}
			case raftpb.MessageType_MsgPropose:
				{

				}
			case raftpb.MessageType_MsgAppend:
				{

				}
			case raftpb.MessageType_MsgAppendResponse:
				{

				}
			case raftpb.MessageType_MsgRequestVote:
				{
					r.mu.Lock()
					if msg.Term > r.currentTerm && msg.Index == r.log.index {
						r.currentTerm = msg.Term
						r.votedFor = msg.From
						//当接收到更大的Term,退位到Follower
						r.becomeFollower()
						go r.peer.SendRaftMessage(msg.From, &raftpb.Message{
							MsgType: raftpb.MessageType_MsgRequestVoteResponse,
							To:      msg.From,
							From:    r.PeerID,
							Term:    r.currentTerm,
							LogTerm: r.log.currentTerm,
							Index:   r.log.index,
							Entries: nil,
							Commit:  r.log.committed,
							Reject:  false,
						})
					} else {
						go r.peer.SendRaftMessage(msg.From, &raftpb.Message{
							MsgType: raftpb.MessageType_MsgRequestVoteResponse,
							To:      msg.From,
							From:    r.PeerID,
							Term:    r.currentTerm,
							LogTerm: r.log.currentTerm,
							Index:   r.log.index,
							Entries: nil,
							Commit:  r.log.committed,
							Reject:  true,
						})
					}
					r.mu.Unlock()
				}
			case raftpb.MessageType_MsgSnapshot:
				{

				}
			case raftpb.MessageType_MsgHeartbeat:
				{

				}
			case raftpb.MessageType_MsgHeartbeatResponse:
				{

				}
			case raftpb.MessageType_MsgTransferLeader:
				{

				}
			case raftpb.MessageType_MsgTimeoutNow:
				{

				}
			}
		}
	}
}

func (r *Raft) HandleCmdMessage(msg *RaftCmdMsg) {

}

func (r *Raft) StartRaft() {

}

// StartElection 发起选举请求
func (r *Raft) StartElection() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.becomeCandidate()
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(r.timeCfg.VoteTimeout)*time.Millisecond)
	res := false
	ctx = context.WithValue(ctx, "res", &res)
	r.es.ctx = ctx
	r.es.cancelFunc = cancelFunc
	for _, peerId := range r.peers {
		go r.peer.SendRaftMessage(peerId, &raftpb.Message{
			MsgType: raftpb.MessageType_MsgRequestVote,
			To:      peerId,
			From:    r.PeerID,
			Term:    r.currentTerm,
			LogTerm: r.log.currentTerm,
			Index:   r.log.committed,
			Entries: nil,
			Commit:  0,
			Reject:  false,
		})
	}
	<-ctx.Done()
	if res {
		r.becomeLeader()
	} else {
		r.becomeFollower()
	}
}

// AppendEntries 处理日志条目追加
func (r *Raft) AppendEntries(term int, entries []LogEntry) {

}

// Candidate发现有比自己更大的任期 or 被大多数拒绝
// Leader发现有比自己更大的任期
func (r *Raft) becomeFollower() {
	r.timer.Enable(utils.ElectionTimerClock)
	r.timer.Disable(utils.HeartBeatClock)
	r.mu.Lock()
	r.State = Follower
	r.mu.Unlock()
	slog.Info("State changed to follower", "PeerId", r.PeerID, "RegionId", r.RegionId)
}

// Term更替：
// 1.Candidate在发起选票过程中，收到更大Term的Leader的心跳信息,退回到Follower
// 2.Follower收到更大Term的心跳信息，且日志完整，于是变更Leader信息
// 3.Leader收到更大Term的心跳信息，且日志完整，于是让位，退回到Follower（由此可见，在每次leader让位的时候，下一个leader的日志总是和上一个leader的日志一致）
func (r *Raft) changeTerm(newTerm uint64, newLeader uint64) {
	r.currentTerm = newTerm
	r.LeaderId = newLeader
}

// 当选举超时的时候从Follower转变为Candidate
// 当当前选举阶段超时的时候开始新的选举阶段
func (r *Raft) becomeCandidate() {
	r.timer.Disable(utils.ElectionTimerClock)
	r.timer.Disable(utils.HeartBeatClock)
	r.mu.Lock()
	r.State = Candidate
	r.votedFor = r.PeerID
	r.currentTerm++
	r.mu.Unlock()
	r.es.setDefault(len(r.peers))
	slog.Info("State changed to candidate", "PeerId", r.PeerID, "RegionId", r.RegionId)
}

// Leader只能从Candidate状态转变
func (r *Raft) becomeLeader() {
	r.timer.Disable(utils.ElectionTimerClock)
	r.timer.Enable(utils.HeartBeatClock)
	r.mu.Lock()
	r.State = Leader
	r.mu.Unlock()
	go r.sendHeartbeat()
	slog.Info("State changed to leader", "PeerId", r.PeerID, "RegionId", r.RegionId)
}

// 领导者发送心跳信号，保持自己的领导者身份
func (r *Raft) sendHeartbeat() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, peerId := range r.peers {
		go r.peer.SendRaftMessage(peerId, &raftpb.Message{
			MsgType: raftpb.MessageType_MsgHeartbeat,
			To:      peerId,
			From:    r.PeerID,
			Term:    r.currentTerm,
			Index:   r.log.index,
			Entries: nil,
			Commit:  0,
			Reject:  false,
		})

	}
}

func (r *Raft) sendAppend() {}

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
