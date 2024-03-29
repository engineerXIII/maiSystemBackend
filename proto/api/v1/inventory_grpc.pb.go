// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.6.1
// source: proto/inventory.proto

package v1

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

const (
	InventoryService_CheckItem_FullMethodName  = "/InventoryService/CheckItem"
	InventoryService_AddItem_FullMethodName    = "/InventoryService/AddItem"
	InventoryService_RemoveItem_FullMethodName = "/InventoryService/RemoveItem"
)

// InventoryServiceClient is the client API for InventoryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InventoryServiceClient interface {
	CheckItem(ctx context.Context, in *ItemRequest, opts ...grpc.CallOption) (*ItemAvailableResponse, error)
	AddItem(ctx context.Context, in *ItemRequest, opts ...grpc.CallOption) (*Response, error)
	RemoveItem(ctx context.Context, in *ItemRequest, opts ...grpc.CallOption) (*Response, error)
}

type inventoryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewInventoryServiceClient(cc grpc.ClientConnInterface) InventoryServiceClient {
	return &inventoryServiceClient{cc}
}

func (c *inventoryServiceClient) CheckItem(ctx context.Context, in *ItemRequest, opts ...grpc.CallOption) (*ItemAvailableResponse, error) {
	out := new(ItemAvailableResponse)
	err := c.cc.Invoke(ctx, InventoryService_CheckItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *inventoryServiceClient) AddItem(ctx context.Context, in *ItemRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, InventoryService_AddItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *inventoryServiceClient) RemoveItem(ctx context.Context, in *ItemRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, InventoryService_RemoveItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InventoryServiceServer is the server API for InventoryService service.
// All implementations must embed UnimplementedInventoryServiceServer
// for forward compatibility
type InventoryServiceServer interface {
	CheckItem(context.Context, *ItemRequest) (*ItemAvailableResponse, error)
	AddItem(context.Context, *ItemRequest) (*Response, error)
	RemoveItem(context.Context, *ItemRequest) (*Response, error)
	mustEmbedUnimplementedInventoryServiceServer()
}

// UnimplementedInventoryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedInventoryServiceServer struct {
}

func (UnimplementedInventoryServiceServer) CheckItem(context.Context, *ItemRequest) (*ItemAvailableResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckItem not implemented")
}
func (UnimplementedInventoryServiceServer) AddItem(context.Context, *ItemRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddItem not implemented")
}
func (UnimplementedInventoryServiceServer) RemoveItem(context.Context, *ItemRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveItem not implemented")
}
func (UnimplementedInventoryServiceServer) mustEmbedUnimplementedInventoryServiceServer() {}

// UnsafeInventoryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InventoryServiceServer will
// result in compilation errors.
type UnsafeInventoryServiceServer interface {
	mustEmbedUnimplementedInventoryServiceServer()
}

func RegisterInventoryServiceServer(s grpc.ServiceRegistrar, srv InventoryServiceServer) {
	s.RegisterService(&InventoryService_ServiceDesc, srv)
}

func _InventoryService_CheckItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServiceServer).CheckItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InventoryService_CheckItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InventoryServiceServer).CheckItem(ctx, req.(*ItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InventoryService_AddItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServiceServer).AddItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InventoryService_AddItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InventoryServiceServer).AddItem(ctx, req.(*ItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InventoryService_RemoveItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServiceServer).RemoveItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InventoryService_RemoveItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InventoryServiceServer).RemoveItem(ctx, req.(*ItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// InventoryService_ServiceDesc is the grpc.ServiceDesc for InventoryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var InventoryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "InventoryService",
	HandlerType: (*InventoryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckItem",
			Handler:    _InventoryService_CheckItem_Handler,
		},
		{
			MethodName: "AddItem",
			Handler:    _InventoryService_AddItem_Handler,
		},
		{
			MethodName: "RemoveItem",
			Handler:    _InventoryService_RemoveItem_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/inventory.proto",
}
