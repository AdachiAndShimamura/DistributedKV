package main

import (
	"AdachiAndShimamura/DistributedKV/kv"
	rpcpb "AdachiAndShimamura/DistributedKV/proto/gen"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	s := grpc.NewServer()
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	service := kv.NewServer()
	rpcpb.RegisterTinyKvRpcServer(s, service)
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
