package rpc

import (
	kvrpc "DistributedKV/proto/gen"
)

type KvServer struct {
	kvrpc.TinyKvServer
}

func (kv *KvServer) Put(key []byte, value []byte) error {}
