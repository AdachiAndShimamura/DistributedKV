package raft

import (
	"context"
	"sync"
)

type electionStatus struct {
	mu         sync.Mutex
	ctx        context.Context
	cancelFunc context.CancelFunc
	agreed     int
	rejected   int
	num        int
	votedMap   map[uint64]bool
}

func (s *electionStatus) statistics(peerID uint64, reject bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if voted := s.votedMap[peerID]; !voted {
		if reject {
			s.rejected++
			if s.rejected >= s.num/2 {
				res := s.ctx.Value("key").(*bool)
				*res = false
				s.cancelFunc()
			}
		} else {
			s.agreed++
			if s.agreed >= s.num/2 {
				res := s.ctx.Value("key").(*bool)
				*res = true
				s.cancelFunc()
			}
		}
	}
}

func (ctx *electionStatus) setDefault(peersNum int) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.agreed = 0
	ctx.rejected = 0
	ctx.num = peersNum
	ctx.votedMap = make(map[uint64]bool, peersNum)
}
