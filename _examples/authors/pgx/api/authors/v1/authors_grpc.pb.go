// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: authors/v1/authors.proto

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
	AuthorsService_CreateAuthor_FullMethodName    = "/authors.v1.AuthorsService/CreateAuthor"
	AuthorsService_DeleteAuthor_FullMethodName    = "/authors.v1.AuthorsService/DeleteAuthor"
	AuthorsService_GetAuthor_FullMethodName       = "/authors.v1.AuthorsService/GetAuthor"
	AuthorsService_ListAuthors_FullMethodName     = "/authors.v1.AuthorsService/ListAuthors"
	AuthorsService_UpdateAuthorBio_FullMethodName = "/authors.v1.AuthorsService/UpdateAuthorBio"
)

// AuthorsServiceClient is the client API for AuthorsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthorsServiceClient interface {
	CreateAuthor(ctx context.Context, in *CreateAuthorRequest, opts ...grpc.CallOption) (*CreateAuthorResponse, error)
	DeleteAuthor(ctx context.Context, in *DeleteAuthorRequest, opts ...grpc.CallOption) (*DeleteAuthorResponse, error)
	GetAuthor(ctx context.Context, in *GetAuthorRequest, opts ...grpc.CallOption) (*GetAuthorResponse, error)
	ListAuthors(ctx context.Context, in *ListAuthorsRequest, opts ...grpc.CallOption) (*ListAuthorsResponse, error)
	UpdateAuthorBio(ctx context.Context, in *UpdateAuthorBioRequest, opts ...grpc.CallOption) (*UpdateAuthorBioResponse, error)
}

type authorsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthorsServiceClient(cc grpc.ClientConnInterface) AuthorsServiceClient {
	return &authorsServiceClient{cc}
}

func (c *authorsServiceClient) CreateAuthor(ctx context.Context, in *CreateAuthorRequest, opts ...grpc.CallOption) (*CreateAuthorResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateAuthorResponse)
	err := c.cc.Invoke(ctx, AuthorsService_CreateAuthor_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorsServiceClient) DeleteAuthor(ctx context.Context, in *DeleteAuthorRequest, opts ...grpc.CallOption) (*DeleteAuthorResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteAuthorResponse)
	err := c.cc.Invoke(ctx, AuthorsService_DeleteAuthor_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorsServiceClient) GetAuthor(ctx context.Context, in *GetAuthorRequest, opts ...grpc.CallOption) (*GetAuthorResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAuthorResponse)
	err := c.cc.Invoke(ctx, AuthorsService_GetAuthor_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorsServiceClient) ListAuthors(ctx context.Context, in *ListAuthorsRequest, opts ...grpc.CallOption) (*ListAuthorsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListAuthorsResponse)
	err := c.cc.Invoke(ctx, AuthorsService_ListAuthors_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authorsServiceClient) UpdateAuthorBio(ctx context.Context, in *UpdateAuthorBioRequest, opts ...grpc.CallOption) (*UpdateAuthorBioResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateAuthorBioResponse)
	err := c.cc.Invoke(ctx, AuthorsService_UpdateAuthorBio_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorsServiceServer is the server API for AuthorsService service.
// All implementations must embed UnimplementedAuthorsServiceServer
// for forward compatibility.
type AuthorsServiceServer interface {
	CreateAuthor(context.Context, *CreateAuthorRequest) (*CreateAuthorResponse, error)
	DeleteAuthor(context.Context, *DeleteAuthorRequest) (*DeleteAuthorResponse, error)
	GetAuthor(context.Context, *GetAuthorRequest) (*GetAuthorResponse, error)
	ListAuthors(context.Context, *ListAuthorsRequest) (*ListAuthorsResponse, error)
	UpdateAuthorBio(context.Context, *UpdateAuthorBioRequest) (*UpdateAuthorBioResponse, error)
	mustEmbedUnimplementedAuthorsServiceServer()
}

// UnimplementedAuthorsServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAuthorsServiceServer struct{}

func (UnimplementedAuthorsServiceServer) CreateAuthor(context.Context, *CreateAuthorRequest) (*CreateAuthorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAuthor not implemented")
}
func (UnimplementedAuthorsServiceServer) DeleteAuthor(context.Context, *DeleteAuthorRequest) (*DeleteAuthorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAuthor not implemented")
}
func (UnimplementedAuthorsServiceServer) GetAuthor(context.Context, *GetAuthorRequest) (*GetAuthorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuthor not implemented")
}
func (UnimplementedAuthorsServiceServer) ListAuthors(context.Context, *ListAuthorsRequest) (*ListAuthorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAuthors not implemented")
}
func (UnimplementedAuthorsServiceServer) UpdateAuthorBio(context.Context, *UpdateAuthorBioRequest) (*UpdateAuthorBioResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAuthorBio not implemented")
}
func (UnimplementedAuthorsServiceServer) mustEmbedUnimplementedAuthorsServiceServer() {}
func (UnimplementedAuthorsServiceServer) testEmbeddedByValue()                        {}

// UnsafeAuthorsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthorsServiceServer will
// result in compilation errors.
type UnsafeAuthorsServiceServer interface {
	mustEmbedUnimplementedAuthorsServiceServer()
}

func RegisterAuthorsServiceServer(s grpc.ServiceRegistrar, srv AuthorsServiceServer) {
	// If the following call pancis, it indicates UnimplementedAuthorsServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AuthorsService_ServiceDesc, srv)
}

func _AuthorsService_CreateAuthor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAuthorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorsServiceServer).CreateAuthor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthorsService_CreateAuthor_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorsServiceServer).CreateAuthor(ctx, req.(*CreateAuthorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthorsService_DeleteAuthor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAuthorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorsServiceServer).DeleteAuthor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthorsService_DeleteAuthor_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorsServiceServer).DeleteAuthor(ctx, req.(*DeleteAuthorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthorsService_GetAuthor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAuthorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorsServiceServer).GetAuthor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthorsService_GetAuthor_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorsServiceServer).GetAuthor(ctx, req.(*GetAuthorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthorsService_ListAuthors_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAuthorsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorsServiceServer).ListAuthors(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthorsService_ListAuthors_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorsServiceServer).ListAuthors(ctx, req.(*ListAuthorsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthorsService_UpdateAuthorBio_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAuthorBioRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorsServiceServer).UpdateAuthorBio(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthorsService_UpdateAuthorBio_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorsServiceServer).UpdateAuthorBio(ctx, req.(*UpdateAuthorBioRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthorsService_ServiceDesc is the grpc.ServiceDesc for AuthorsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthorsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "authors.v1.AuthorsService",
	HandlerType: (*AuthorsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAuthor",
			Handler:    _AuthorsService_CreateAuthor_Handler,
		},
		{
			MethodName: "DeleteAuthor",
			Handler:    _AuthorsService_DeleteAuthor_Handler,
		},
		{
			MethodName: "GetAuthor",
			Handler:    _AuthorsService_GetAuthor_Handler,
		},
		{
			MethodName: "ListAuthors",
			Handler:    _AuthorsService_ListAuthors_Handler,
		},
		{
			MethodName: "UpdateAuthorBio",
			Handler:    _AuthorsService_UpdateAuthorBio_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "authors/v1/authors.proto",
}
