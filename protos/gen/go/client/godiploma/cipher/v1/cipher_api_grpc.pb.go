// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: client/godiploma/cipher/v1/cipher_api.proto

package cipherv1

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
	CipherService_CreateStegoImage_FullMethodName = "/client.godiploma.cipher.v1.CipherService/CreateStegoImage"
	CipherService_Extract_FullMethodName          = "/client.godiploma.cipher.v1.CipherService/Extract"
)

// CipherServiceClient is the client API for CipherService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// CipherService steganography cipher service.
type CipherServiceClient interface {
	// CreateStegoImage creates a stego image.
	CreateStegoImage(ctx context.Context, in *CreateStegoImageRequest, opts ...grpc.CallOption) (*CreateStegoImageResponse, error)
	// Extract extracts plaintext from a stego image.
	Extract(ctx context.Context, in *ExtractRequest, opts ...grpc.CallOption) (*ExtractResponse, error)
}

type cipherServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCipherServiceClient(cc grpc.ClientConnInterface) CipherServiceClient {
	return &cipherServiceClient{cc}
}

func (c *cipherServiceClient) CreateStegoImage(ctx context.Context, in *CreateStegoImageRequest, opts ...grpc.CallOption) (*CreateStegoImageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateStegoImageResponse)
	err := c.cc.Invoke(ctx, CipherService_CreateStegoImage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cipherServiceClient) Extract(ctx context.Context, in *ExtractRequest, opts ...grpc.CallOption) (*ExtractResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ExtractResponse)
	err := c.cc.Invoke(ctx, CipherService_Extract_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CipherServiceServer is the server API for CipherService service.
// All implementations should embed UnimplementedCipherServiceServer
// for forward compatibility.
//
// CipherService steganography cipher service.
type CipherServiceServer interface {
	// CreateStegoImage creates a stego image.
	CreateStegoImage(context.Context, *CreateStegoImageRequest) (*CreateStegoImageResponse, error)
	// Extract extracts plaintext from a stego image.
	Extract(context.Context, *ExtractRequest) (*ExtractResponse, error)
}

// UnimplementedCipherServiceServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCipherServiceServer struct{}

func (UnimplementedCipherServiceServer) CreateStegoImage(context.Context, *CreateStegoImageRequest) (*CreateStegoImageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateStegoImage not implemented")
}
func (UnimplementedCipherServiceServer) Extract(context.Context, *ExtractRequest) (*ExtractResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Extract not implemented")
}
func (UnimplementedCipherServiceServer) testEmbeddedByValue() {}

// UnsafeCipherServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CipherServiceServer will
// result in compilation errors.
type UnsafeCipherServiceServer interface {
	mustEmbedUnimplementedCipherServiceServer()
}

func RegisterCipherServiceServer(s grpc.ServiceRegistrar, srv CipherServiceServer) {
	// If the following call pancis, it indicates UnimplementedCipherServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CipherService_ServiceDesc, srv)
}

func _CipherService_CreateStegoImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateStegoImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CipherServiceServer).CreateStegoImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CipherService_CreateStegoImage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CipherServiceServer).CreateStegoImage(ctx, req.(*CreateStegoImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CipherService_Extract_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExtractRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CipherServiceServer).Extract(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CipherService_Extract_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CipherServiceServer).Extract(ctx, req.(*ExtractRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CipherService_ServiceDesc is the grpc.ServiceDesc for CipherService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CipherService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "client.godiploma.cipher.v1.CipherService",
	HandlerType: (*CipherServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateStegoImage",
			Handler:    _CipherService_CreateStegoImage_Handler,
		},
		{
			MethodName: "Extract",
			Handler:    _CipherService_Extract_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "client/godiploma/cipher/v1/cipher_api.proto",
}
