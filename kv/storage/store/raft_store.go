package store

import (
	"github.com/robfig/cron/v3"
)

// RaftStore 表示作为存储节点的一台机器，以StoreID标识每台机器
type RaftStore struct {
	storeID uint64
	router  *Router
	timer   *cron.Cron
}

func (r *RaftStore) Start() {
	timer := cron.New()

}

func (r *RaftStore) AddTimer() {
	
}

func (r *RaftStore) Stop() {
	r.timer.Stop()
}
