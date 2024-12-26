package rpc

import (
	rpcpb "AdachiAndShimamura/DistributedKV/proto/gen"
	"AdachiAndShimamura/DistributedKV/proto/gen/raftpb"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"sync"
	"time"
)

type RaftConn struct {
	streamMu sync.Mutex
	stream   rpcpb.TinyKvRpc_RaftClient
	ctx      context.Context
	cancel   context.CancelFunc
}

func newRaftConn(addr string) (*RaftConn, error) {
	cc, err := grpc.NewClient(addr,
		//禁用TLS
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                3 * time.Second,
			Timeout:             60 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := rpcpb.NewTinyKvRpcClient(cc).Raft(ctx)
	if err != nil {
		cancel()
		return nil, err
	}
	return &RaftConn{
		streamMu: sync.Mutex{},
		stream:   stream,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}
func (c *RaftConn) Send(message *raftpb.RaftMessage) error {
	c.streamMu.Lock()
	defer c.streamMu.Unlock()
	return c.stream.Send(message)
}

func (c *RaftConn) Stop() {
	c.cancel()
}

type RaftClient struct {
	mu sync.RWMutex
	//storeID->RPC-Client
	conns map[uint64]*RaftConn
	//StoreID->Addr
	addrMap map[uint64]string
}

// GetConn 获取连接，如果没有，则创建连接
func (c *RaftClient) GetConn(id uint64) (*RaftConn, error) {
	c.mu.RLock()
	conn, ok := c.conns[id]
	if ok {
		return conn, nil
	}
	addr, ok := c.addrMap[id]
	if !ok {
		return nil, fmt.Errorf("id:%d not exist", id)
	}
	c.mu.RUnlock()
	newConn, err := newRaftConn(addr)
	if err != nil {
		return nil, err
	}
	//再次检查，防止多线程环境下多次创建连接（因为创建连接耗时较高，因此避免在创建连接时加锁）
	c.mu.Lock()
	defer c.mu.Unlock()
	if conn, ok := c.conns[id]; ok {
		newConn.Stop()
		return conn, nil
	}
	c.conns[id] = conn
	return conn, nil
}

func (c *RaftClient) Send(id uint64, message *raftpb.RaftMessage) error {
	conn, err := c.GetConn(id)
	if err == nil {
		return conn.Send(message)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.conns[id]; ok {
		delete(c.conns, id)
	}
	return err
}

func (c *RaftClient) GetAddr(id uint64) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	addr, ok := c.addrMap[id]
	if ok {
		return addr, nil
	}
	return "", errors.New("id:%d not exist")
}

func (c *RaftClient) AddAddr(id uint64, addr string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.addrMap[id] = addr
}
