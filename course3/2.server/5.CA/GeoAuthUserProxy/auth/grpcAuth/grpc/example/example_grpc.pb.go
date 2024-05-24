// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package __

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

// TokenValidationServiceClient is the client API for TokenValidationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TokenValidationServiceClient interface {
	ValidateToken(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*TokenResponse, error)
}

type tokenValidationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTokenValidationServiceClient(cc grpc.ClientConnInterface) TokenValidationServiceClient {
	return &tokenValidationServiceClient{cc}
}

func (c *tokenValidationServiceClient) ValidateToken(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*TokenResponse, error) {
	out := new(TokenResponse)
	err := c.cc.Invoke(ctx, "/main.TokenValidationService/ValidateToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TokenValidationServiceServer is the server API for TokenValidationService service.
// All implementations must embed UnimplementedTokenValidationServiceServer
// for forward compatibility
type TokenValidationServiceServer interface {
	ValidateToken(context.Context, *TokenRequest) (*TokenResponse, error)
	mustEmbedUnimplementedTokenValidationServiceServer()
}

// UnimplementedTokenValidationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTokenValidationServiceServer struct {
}

func (UnimplementedTokenValidationServiceServer) ValidateToken(context.Context, *TokenRequest) (*TokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateToken not implemented")
}
func (UnimplementedTokenValidationServiceServer) mustEmbedUnimplementedTokenValidationServiceServer() {
}

// UnsafeTokenValidationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TokenValidationServiceServer will
// result in compilation errors.
type UnsafeTokenValidationServiceServer interface {
	mustEmbedUnimplementedTokenValidationServiceServer()
}

func RegisterTokenValidationServiceServer(s grpc.ServiceRegistrar, srv TokenValidationServiceServer) {
	s.RegisterService(&TokenValidationService_ServiceDesc, srv)
}

func _TokenValidationService_ValidateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenValidationServiceServer).ValidateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.TokenValidationService/ValidateToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenValidationServiceServer).ValidateToken(ctx, req.(*TokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TokenValidationService_ServiceDesc is the grpc.ServiceDesc for TokenValidationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TokenValidationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.TokenValidationService",
	HandlerType: (*TokenValidationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ValidateToken",
			Handler:    _TokenValidationService_ValidateToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "example.proto",
}