// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: api/ledger/v1/service.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	LedgerService_Transact_FullMethodName = "/go_example.api.ledger.v1.LedgerService/Transact"
)

// LedgerServiceClient is the client API for LedgerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LedgerServiceClient interface {
	Transact(ctx context.Context, in *TransactRequest, opts ...grpc.CallOption) (*TransactResponse, error)
}

type ledgerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLedgerServiceClient(cc grpc.ClientConnInterface) LedgerServiceClient {
	return &ledgerServiceClient{cc}
}

func (c *ledgerServiceClient) Transact(ctx context.Context, in *TransactRequest, opts ...grpc.CallOption) (*TransactResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TransactResponse)
	err := c.cc.Invoke(ctx, LedgerService_Transact_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LedgerServiceServer is the server API for LedgerService service.
// All implementations must embed UnimplementedLedgerServiceServer
// for forward compatibility.
type LedgerServiceServer interface {
	Transact(context.Context, *TransactRequest) (*TransactResponse, error)
	mustEmbedUnimplementedLedgerServiceServer()
}

// UnimplementedLedgerServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLedgerServiceServer struct{}

func (UnimplementedLedgerServiceServer) Transact(context.Context, *TransactRequest) (*TransactResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Transact not implemented")
}
func (UnimplementedLedgerServiceServer) mustEmbedUnimplementedLedgerServiceServer() {}
func (UnimplementedLedgerServiceServer) testEmbeddedByValue()                       {}

// UnsafeLedgerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LedgerServiceServer will
// result in compilation errors.
type UnsafeLedgerServiceServer interface {
	mustEmbedUnimplementedLedgerServiceServer()
}

func RegisterLedgerServiceServer(s grpc.ServiceRegistrar, srv LedgerServiceServer) {
	// If the following call pancis, it indicates UnimplementedLedgerServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&LedgerService_ServiceDesc, srv)
}

func _LedgerService_Transact_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransactRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServiceServer).Transact(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LedgerService_Transact_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServiceServer).Transact(ctx, req.(*TransactRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LedgerService_ServiceDesc is the grpc.ServiceDesc for LedgerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LedgerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "go_example.api.ledger.v1.LedgerService",
	HandlerType: (*LedgerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Transact",
			Handler:    _LedgerService_Transact_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/ledger/v1/service.proto",
}