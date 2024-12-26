package main

import (
	rpcpb "AdachiAndShimamura/DistributedKV/proto/gen"
	"AdachiAndShimamura/DistributedKV/proto/gen/kvrpcpb"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type KvClient struct {
	kvs []rpcpb.TinyKvRpcClient
}

func BuildClient() *KvClient {
	return &KvClient{kvs: make([]rpcpb.TinyKvRpcClient, 0)}
}
func (c *KvClient) Connect() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("127.0.0.1:8080", opts...)
	if err != nil {
		log.Fatalln("error:", err)
	}
	client := rpcpb.NewTinyKvRpcClient(conn)
	c.kvs = append(c.kvs, client)
	//defer func(client *grpc.ClientConn) {
	//	err := client.Close()
	//	if err != nil {
	//		log.Print("error:", err)
	//	}
	//}(conn)
}

func main() {
	client := BuildClient()
	client.Connect()
	cxt := context.Background()
	client.kvs[0].RawPut(cxt, &kvrpcpb.RawPut{
		Key:   []byte("111"),
		Value: []byte("222"),
	})
	res, _ := client.kvs[0].RawGet(cxt, &kvrpcpb.RawGet{
		Key: []byte("111"),
	})
	println(string(res.Value))
}
