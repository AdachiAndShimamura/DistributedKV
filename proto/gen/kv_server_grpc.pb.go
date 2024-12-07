// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.29.0--rc2
// source: kv_server.proto

package kv_rpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TinyKvClient is the client API for TinyKv service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TinyKvClient interface {
	// RawKV commands
	RawGet(ctx context.Context, in *RawGet, opts ...grpc.CallOption) (*RawGetResponse, error)
	RawPut(ctx context.Context, in *RawPut, opts ...grpc.CallOption) (*RawPutResponse, error)
	RawDelete(ctx context.Context, in *RawDelete, opts ...grpc.CallOption) (*RawDeleteResponse, error)
	RawScan(ctx context.Context, in *RawScan, opts ...grpc.CallOption) (*RawScanResponse, error)
}

type tinyKvClient struct {
	cc grpc.ClientConnInterface
}

func NewTinyKvClient(cc grpc.ClientConnInterface) TinyKvClient {
	return &tinyKvClient{cc}
}

func (c *tinyKvClient) RawGet(ctx context.Context, in *RawGet, opts ...grpc.CallOption) (*RawGetResponse, error) {
	out := new(RawGetResponse)
	err := c.cc.Invoke(ctx, "/kv_rpc.TinyKv/RawGet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tinyKvClient) RawPut(ctx context.Context, in *RawPut, opts ...grpc.CallOption) (*RawPutResponse, error) {
	out := new(RawPutResponse)
	err := c.cc.Invoke(ctx, "/kv_rpc.TinyKv/RawPut", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tinyKvClient) RawDelete(ctx context.Context, in *RawDelete, opts ...grpc.CallOption) (*RawDeleteResponse, error) {
	out := new(RawDeleteResponse)
	err := c.cc.Invoke(ctx, "/kv_rpc.TinyKv/RawDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tinyKvClient) RawScan(ctx context.Context, in *RawScan, opts ...grpc.CallOption) (*RawScanResponse, error) {
	out := new(RawScanResponse)
	err := c.cc.Invoke(ctx, "/kv_rpc.TinyKv/RawScan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TinyKvServer is the server API for TinyKv service.
// All implementations must embed UnimplementedTinyKvServer
// for forward compatibility
type TinyKvServer interface {
	// RawKV commands
	RawGet(context.Context, *RawGet) (*RawGetResponse, error)
	RawPut(context.Context, *RawPut) (*RawPutResponse, error)
	RawDelete(context.Context, *RawDelete) (*RawDeleteResponse, error)
	RawScan(context.Context, *RawScan) (*RawScanResponse, error)
	mustEmbedUnimplementedTinyKvServer()
}

// UnimplementedTinyKvServer must be embedded to have forward compatible implementations.
type UnimplementedTinyKvServer struct {
}

func (UnimplementedTinyKvServer) RawGet(context.Context, *RawGet) (*RawGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RawGet not implemented")
}
func (UnimplementedTinyKvServer) RawPut(context.Context, *RawPut) (*RawPutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RawPut not implemented")
}
func (UnimplementedTinyKvServer) RawDelete(context.Context, *RawDelete) (*RawDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RawDelete not implemented")
}
func (UnimplementedTinyKvServer) RawScan(context.Context, *RawScan) (*RawScanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RawScan not implemented")
}
func (UnimplementedTinyKvServer) mustEmbedUnimplementedTinyKvServer() {}

// UnsafeTinyKvServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TinyKvServer will
// result in compilation errors.
type UnsafeTinyKvServer interface {
	mustEmbedUnimplementedTinyKvServer()
}

func RegisterTinyKvServer(s grpc.ServiceRegistrar, srv TinyKvServer) {
	s.RegisterService(&TinyKv_ServiceDesc, srv)
}

func _TinyKv_RawGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RawGet)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TinyKvServer).RawGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kv_rpc.TinyKv/RawGet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TinyKvServer).RawGet(ctx, req.(*RawGet))
	}
	return interceptor(ctx, in, info, handler)
}

func _TinyKv_RawPut_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RawPut)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TinyKvServer).RawPut(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kv_rpc.TinyKv/RawPut",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TinyKvServer).RawPut(ctx, req.(*RawPut))
	}
	return interceptor(ctx, in, info, handler)
}

func _TinyKv_RawDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RawDelete)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TinyKvServer).RawDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kv_rpc.TinyKv/RawDelete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TinyKvServer).RawDelete(ctx, req.(*RawDelete))
	}
	return interceptor(ctx, in, info, handler)
}

func _TinyKv_RawScan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RawScan)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TinyKvServer).RawScan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kv_rpc.TinyKv/RawScan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TinyKvServer).RawScan(ctx, req.(*RawScan))
	}
	return interceptor(ctx, in, info, handler)
}

// TinyKv_ServiceDesc is the grpc.ServiceDesc for TinyKv service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TinyKv_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kv_rpc.TinyKv",
	HandlerType: (*TinyKvServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RawGet",
			Handler:    _TinyKv_RawGet_Handler,
		},
		{
			MethodName: "RawPut",
			Handler:    _TinyKv_RawPut_Handler,
		},
		{
			MethodName: "RawDelete",
			Handler:    _TinyKv_RawDelete_Handler,
		},
		{
			MethodName: "RawScan",
			Handler:    _TinyKv_RawScan_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kv_server.proto",
}
