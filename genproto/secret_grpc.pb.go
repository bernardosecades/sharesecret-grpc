// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

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

// SecretServiceClient is the client API for SecretService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SecretServiceClient interface {
	CreateSecret(ctx context.Context, in *CreateSecretRequest, opts ...grpc.CallOption) (*CreateSecretResponse, error)
	SeeSecret(ctx context.Context, in *SeeSecretRequest, opts ...grpc.CallOption) (*SeeSecretResponse, error)
}

type secretServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSecretServiceClient(cc grpc.ClientConnInterface) SecretServiceClient {
	return &secretServiceClient{cc}
}

func (c *secretServiceClient) CreateSecret(ctx context.Context, in *CreateSecretRequest, opts ...grpc.CallOption) (*CreateSecretResponse, error) {
	out := new(CreateSecretResponse)
	err := c.cc.Invoke(ctx, "/sharesecret.SecretService/CreateSecret", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretServiceClient) SeeSecret(ctx context.Context, in *SeeSecretRequest, opts ...grpc.CallOption) (*SeeSecretResponse, error) {
	out := new(SeeSecretResponse)
	err := c.cc.Invoke(ctx, "/sharesecret.SecretService/SeeSecret", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SecretServiceServer is the server API for SecretService service.
// All implementations should embed UnimplementedSecretServiceServer
// for forward compatibility
type SecretServiceServer interface {
	CreateSecret(context.Context, *CreateSecretRequest) (*CreateSecretResponse, error)
	SeeSecret(context.Context, *SeeSecretRequest) (*SeeSecretResponse, error)
}

// UnimplementedSecretServiceServer should be embedded to have forward compatible implementations.
type UnimplementedSecretServiceServer struct {
}

func (UnimplementedSecretServiceServer) CreateSecret(context.Context, *CreateSecretRequest) (*CreateSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSecret not implemented")
}
func (UnimplementedSecretServiceServer) SeeSecret(context.Context, *SeeSecretRequest) (*SeeSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SeeSecret not implemented")
}

// UnsafeSecretServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SecretServiceServer will
// result in compilation errors.
type UnsafeSecretServiceServer interface {
	mustEmbedUnimplementedSecretServiceServer()
}

func RegisterSecretServiceServer(s grpc.ServiceRegistrar, srv SecretServiceServer) {
	s.RegisterService(&SecretService_ServiceDesc, srv)
}

func _SecretService_CreateSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretServiceServer).CreateSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sharesecret.SecretService/CreateSecret",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretServiceServer).CreateSecret(ctx, req.(*CreateSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SecretService_SeeSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SeeSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretServiceServer).SeeSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sharesecret.SecretService/SeeSecret",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretServiceServer).SeeSecret(ctx, req.(*SeeSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SecretService_ServiceDesc is the grpc.ServiceDesc for SecretService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SecretService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sharesecret.SecretService",
	HandlerType: (*SecretServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateSecret",
			Handler:    _SecretService_CreateSecret_Handler,
		},
		{
			MethodName: "SeeSecret",
			Handler:    _SecretService_SeeSecret_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "secret.proto",
}
