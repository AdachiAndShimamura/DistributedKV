package storage

import (
	"AdachiAndShimamura/DistributedKV/kv/engine_util"
	"AdachiAndShimamura/DistributedKV/kv/storage/store"
	rpcpb "AdachiAndShimamura/DistributedKV/proto/gen"
	"AdachiAndShimamura/DistributedKV/proto/gen/kvrpcpb"
	"AdachiAndShimamura/DistributedKV/proto/gen/raftpb"
	"io"
	"os"
	"path/filepath"
)

// RaftStorage 单机存储的Raft版本
type RaftStorage struct {
	//实际存储引擎，包含一个Kv存储引擎和一个Raft元数据存储引擎
	storage *engine_util.Engines
	config  *Config
	store   *store.RaftStore
	router  *store.Router
}

func (s *RaftStorage) Start() error {

}

func (s *RaftStorage) Stop() error {
	//TODO implement me
	panic("implement me")
}

// 存储引擎的写接口，会阻塞等待raft peer处理
func (s *RaftStorage) Write(ctx *kvrpcpb.Context, batch []Modify) error {
	var cmds []*raftpb.Request
	for _, cmd := range batch {
		switch cmd.Data.(type) {
		case Put:
			{
				put := cmd.Data.(Put)
				cmds = append(cmds, &raftpb.Request{
					CmdType: raftpb.RaftCmdType_Put,
					Put: &raftpb.PutReq{
						Key:   put.Key,
						Value: put.Value,
						Cf:    put.Cf,
					},
				})
			}
		case Delete:
			{
				del := cmd.Data.(Delete)
				cmds = append(cmds, &raftpb.Request{
					CmdType: raftpb.RaftCmdType_Delete,
					Del: &raftpb.DeleteReq{
						Key: del.Key,
						Cf:  del.Cf,
					},
				})
			}
		}
	}
	header := &raftpb.RequestHeader{RegionId: ctx.RegionId}
	req := &raftpb.RaftCmdRequest{
		Header:   header,
		Requests: cmds,
	}
	cb := store.NewCallBack()
	err := s.HandleCmdMessage(req, cb)
	if err != nil {
		return err
	}
	return store.CheckCmdResp(cb.WaitResp(), len(cmds))
}

// Reader 存储引擎的读接口，会阻塞等待raft peer进行状态检查
func (s *RaftStorage) Reader(ctx *kvrpcpb.Context) (StorageReader, error) {
	//TODO implement me
	panic("implement me")
}

func NewRaftStorage(config *Config) *RaftStorage {
	dbPath := config.DBPath
	raftPath := filepath.Join(dbPath, "raft")
	kvPath := filepath.Join(dbPath, "kv")
	os.Mkdir(raftPath, 0700)
	os.Mkdir(kvPath, 0700)
	storage := engine_util.NewEngines(raftPath, kvPath)
	return &RaftStorage{
		storage: storage,
		config:  config,
	}
}

func (s *RaftStorage) HandleRaftMessage(stream *rpcpb.TinyKvRpc_RaftServer) error {
	for {
		data, err := (*stream).Recv()
		if err == io.EOF {
			err = (*stream).SendAndClose(&raftpb.None{})
		}
		if err != nil {
			return err
		}
		s.router.SendRaftMessage(data)
	}
}

func (s *RaftStorage) HandleCmdMessage(cmd *raftpb.RaftCmdRequest, cb *store.CallBack) error {
	return s.router.SendRaftCmdMessage(cmd, cb)
}
