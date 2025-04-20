package utils

import (
	"AdachiAndShimamura/DistributedKV/kv/common"
	"math/rand/v2"
	"sync"
	"time"
)

var TickerManager TickerCore

type PeerTicker struct {
	peerID uint64
	//时间间隔的基值，ms为单位
	baseInterval uint64
	clocks       map[ClockType]*clock
	ch           chan struct{}
	peerCh       chan ClockType
}
type ClockType uint8

const (
	HeartBeatClock ClockType = iota
	ElectionTimerClock
)

func NewPeerTicker(peerID uint64, cfg *common.ClockConfig, ch chan ClockType) *PeerTicker {
	clocks := make(map[ClockType]*clock, 3)
	//心跳信息发送间隔
	clocks[HeartBeatClock] = &clock{
		runAt:  cfg.HeartbeatInterval / cfg.BaseTimeInternal,
		enable: false,
		fn: func() {
			ch <- HeartBeatClock
		},
	}
	//选举超时时间
	clocks[ElectionTimerClock] = &clock{
		runAt:  rand.N((cfg.ElectionTimeoutEnd-cfg.ElectionTimeoutStart)/cfg.BaseTimeInternal) + cfg.ElectionTimeoutStart/cfg.BaseTimeInternal,
		enable: false,
		fn: func() {
			ch <- ElectionTimerClock
		},
	}
	return &PeerTicker{
		peerID:       peerID,
		baseInterval: cfg.BaseTimeInternal,
		clocks:       clocks,
	}
}

type clock struct {
	runAt  uint64
	now    uint64
	enable bool
	fn     func()
}

func (t *PeerTicker) Start() {
	ch := make(chan struct{}, 1)
	TickerManager.RegisterPeer(t.peerID, ch)
	for {
		<-t.ch
		for _, clock := range t.clocks {
			if !clock.enable {
				continue
			}
			clock.Tick()
		}
	}
}
func (t *PeerTicker) Stop() {
	TickerManager.Unregister(t.peerID)
}
func (t *PeerTicker) SetFn(tickerType ClockType, fn func()) {
	t.clocks[tickerType].fn = fn
}
func (t *clock) Tick() {
	t.now++
	if t.runAt == t.now {
		t.now = 0
		t.fn()
	}
}

func (t *PeerTicker) Enable(tickerType ClockType) {
	t.clocks[tickerType].enable = true
}

func (t *PeerTicker) Disable(tickerType ClockType) {
	t.clocks[tickerType].now = 0
	t.clocks[tickerType].enable = false
}

func (t *PeerTicker) Reset(tickerType ClockType) {
	t.clocks[tickerType].now = 0
}

type TickerCore struct {
	senders      sync.Map
	baseInterval uint64
	ticker       *time.Ticker
}

func (t *TickerCore) Start(cfg *common.ClockConfig) {
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
