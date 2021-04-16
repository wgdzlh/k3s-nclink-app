// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package configmodel

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

// ModelDistClient is the client API for ModelDist service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ModelDistClient interface {
	GetModel(ctx context.Context, in *ModelRequest, opts ...grpc.CallOption) (*ModelReply, error)
}

type modelDistClient struct {
	cc grpc.ClientConnInterface
}

func NewModelDistClient(cc grpc.ClientConnInterface) ModelDistClient {
	return &modelDistClient{cc}
}

func (c *modelDistClient) GetModel(ctx context.Context, in *ModelRequest, opts ...grpc.CallOption) (*ModelReply, error) {
	out := new(ModelReply)
	err := c.cc.Invoke(ctx, "/configmodel.ModelDist/GetModel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ModelDistServer is the server API for ModelDist service.
// All implementations must embed UnimplementedModelDistServer
// for forward compatibility
type ModelDistServer interface {
	GetModel(context.Context, *ModelRequest) (*ModelReply, error)
	mustEmbedUnimplementedModelDistServer()
}

// UnimplementedModelDistServer must be embedded to have forward compatible implementations.
type UnimplementedModelDistServer struct {
}

func (UnimplementedModelDistServer) GetModel(context.Context, *ModelRequest) (*ModelReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetModel not implemented")
}
func (UnimplementedModelDistServer) mustEmbedUnimplementedModelDistServer() {}

// UnsafeModelDistServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ModelDistServer will
// result in compilation errors.
type UnsafeModelDistServer interface {
	mustEmbedUnimplementedModelDistServer()
}

func RegisterModelDistServer(s grpc.ServiceRegistrar, srv ModelDistServer) {
	s.RegisterService(&ModelDist_ServiceDesc, srv)
}

func _ModelDist_GetModel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ModelDistServer).GetModel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/configmodel.ModelDist/GetModel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ModelDistServer).GetModel(ctx, req.(*ModelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ModelDist_ServiceDesc is the grpc.ServiceDesc for ModelDist service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ModelDist_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "configmodel.ModelDist",
	HandlerType: (*ModelDistServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetModel",
			Handler:    _ModelDist_GetModel_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "configmodel-dist.proto",
}
