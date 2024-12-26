package utils

import (
	"AdachiAndShimamura/DistributedKV/kv/storage"
	"sync"
	"time"
)

var TickerManager TickerCore

type PeerTicker struct {
	peerID uint64
	//时间间隔的基值，ms为单位
	baseInterval uint64
	schedulers   []*tickScheduler
	ch           chan struct{}
}
type TickType uint8

const (
	HeartBeat TickType = iota
	electionTimer
)

func NewPeerTicker(peerID uint64, baseInterval uint64) *PeerTicker {
	return &PeerTicker{
		peerID:       peerID,
		baseInterval: baseInterval,
		schedulers:   make([]*tickScheduler, 2),
	}
}

type tickScheduler struct {
	runAt uint64
	now   uint64
	fn    func()
}

func (t *PeerTicker) Start() {
	ch := make(chan struct{}, 1)
	TickerManager.RegisterPeer(t.peerID, ch)

	for _, scheduler := range t.schedulers {
		if scheduler != nil {
			scheduler.Tick()
		}
	}
}
func (t *PeerTicker) Stop() {
	TickerManager.Unregister(t.peerID)
}

func (t *tickScheduler) Tick() {
	t.now++
	if t.runAt == t.now {
		t.now = 0
		go t.fn()
	}
}

func (t *PeerTicker) Register(msgType TickType, time uint64, fn func()) {
	t.schedulers[msgType] = &tickScheduler{
		runAt: time / t.baseInterval,
		now:   0,
		fn:    fn,
	}
}

func (t *PeerTicker) Unregister(msgType TickType) {
	t.schedulers[msgType] = nil
}

type TickerCore struct {
	senders      sync.Map
	baseInterval uint64
	ticker       *time.Ticker
}

func (t *TickerCore) Start(cfg storage.Config) {
	t.baseInterval = cfg.BaseTimeInternal
	t.ticker = time.NewTicker(time.Duration(t.baseInterval) * time.Millisecond)
	ch := t.ticker.C
	for {
		<-ch
		t.senders.Range(func(_, value interface{}) bool {
			value.(chan struct{}) <- struct{}{}
			return true
		})

	}
}

func (t *TickerCore) RegisterPeer(peerID uint64, ch chan struct{}) {
	t.senders.Store(peerID, ch)
}

func (t *TickerCore) Unregister(peerID uint64) {
	t.senders.Delete(peerID)
}
