// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.1.0
// - protoc             v3.14.0
// source: client.proto

package chickaree

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

// ChickareeDBClient is the client API for ChickareeDB service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChickareeDBClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error)
}

type chickareeDBClient struct {
	cc grpc.ClientConnInterface
}

func NewChickareeDBClient(cc grpc.ClientConnInterface) ChickareeDBClient {
	return &chickareeDBClient{cc}
}

func (c *chickareeDBClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/client.v1.ChickareeDB/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chickareeDBClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error) {
	out := new(SetResponse)
	err := c.cc.Invoke(ctx, "/client.v1.ChickareeDB/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChickareeDBServer is the server API for ChickareeDB service.
// All implementations must embed UnimplementedChickareeDBServer
// for forward compatibility
type ChickareeDBServer interface {
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Set(context.Context, *SetRequest) (*SetResponse, error)
	mustEmbedUnimplementedChickareeDBServer()
}

// UnimplementedChickareeDBServer must be embedded to have forward compatible implementations.
type UnimplementedChickareeDBServer struct {
}

func (UnimplementedChickareeDBServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedChickareeDBServer) Set(context.Context, *SetRequest) (*SetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedChickareeDBServer) mustEmbedUnimplementedChickareeDBServer() {}

// UnsafeChickareeDBServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChickareeDBServer will
// result in compilation errors.
type UnsafeChickareeDBServer interface {
	mustEmbedUnimplementedChickareeDBServer()
}

func RegisterChickareeDBServer(s grpc.ServiceRegistrar, srv ChickareeDBServer) {
	s.RegisterService(&ChickareeDB_ServiceDesc, srv)
}

func _ChickareeDB_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChickareeDBServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/client.v1.ChickareeDB/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChickareeDBServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChickareeDB_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChickareeDBServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/client.v1.ChickareeDB/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChickareeDBServer).Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChickareeDB_ServiceDesc is the grpc.ServiceDesc for ChickareeDB service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChickareeDB_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "client.v1.ChickareeDB",
	HandlerType: (*ChickareeDBServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _ChickareeDB_Get_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _ChickareeDB_Set_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "client.proto",
}