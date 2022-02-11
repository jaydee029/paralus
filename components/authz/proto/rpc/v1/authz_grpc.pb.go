// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package authzv1

import (
	context "context"
	types "github.com/RafaySystems/rcloud-base/components/authz/proto/types"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AuthzClient is the client API for Authz service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthzClient interface {
	Enforce(ctx context.Context, in *types.EnforceRequest, opts ...grpc.CallOption) (*types.BoolReply, error)
	// List Policies accpets Policy whose fileds are used for filtering
	// Filtering is done per field for the policy
	// For Example:
	// The Policy obj:
	//    sub => ""
	//    ns => ""
	//    proj => project1
	//    org => org1
	//    obj => ""
	//    act => ""
	// Returns policies related to project1 and org1 (Empty string matches all)
	ListPolicies(ctx context.Context, in *types.Policy, opts ...grpc.CallOption) (*types.Policies, error)
	CreatePolicies(ctx context.Context, in *types.Policies, opts ...grpc.CallOption) (*types.BoolReply, error)
	DeletePolicies(ctx context.Context, in *types.Policy, opts ...grpc.CallOption) (*types.BoolReply, error)
	ListUserGroups(ctx context.Context, in *types.UserGroup, opts ...grpc.CallOption) (*types.UserGroups, error)
	CreateUserGroups(ctx context.Context, in *types.UserGroups, opts ...grpc.CallOption) (*types.BoolReply, error)
	DeleteUserGroups(ctx context.Context, in *types.UserGroup, opts ...grpc.CallOption) (*types.BoolReply, error)
	ListRolePermissionMappings(ctx context.Context, in *types.FilteredRolePermissionMapping, opts ...grpc.CallOption) (*types.RolePermissionMappingList, error)
	CreateRolePermissionMappings(ctx context.Context, in *types.RolePermissionMappingList, opts ...grpc.CallOption) (*types.BoolReply, error)
	DeleteRolePermissionMappings(ctx context.Context, in *types.FilteredRolePermissionMapping, opts ...grpc.CallOption) (*types.BoolReply, error)
}

type authzClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthzClient(cc grpc.ClientConnInterface) AuthzClient {
	return &authzClient{cc}
}

func (c *authzClient) Enforce(ctx context.Context, in *types.EnforceRequest, opts ...grpc.CallOption) (*types.BoolReply, error) {
	out := new(types.BoolReply)
	err := c.cc.Invoke(ctx, "/rafay.dev.rpc.authz.v1.Authz/Enforce", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authzClient) ListPolicies(ctx context.Context, in *types.Policy, opts ...grpc.CallOption) (*types.Policies, error) {
	out := new(types.Policies)
	err := c.cc.Invoke(ctx, "/rafay.dev.rpc.authz.v1.Authz/ListPolicies", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authzClient) CreatePolicies(ctx context.Context, in *types.Policies, opts ...grpc.CallOption) (*types.BoolReply, error) {
	out := new(types.BoolReply)
	err := c.cc.Invoke(ctx, "/rafay.dev.rpc.authz.v1.Authz/CreatePolicies", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authzClient) DeletePolicies(ctx context.Context, in *types.Policy, opts ...grpc.CallOption) (*types.BoolReply, error) {
	out := new(types.BoolReply)
	err := c.cc.Invoke(ctx, "/rafay.dev.rpc.authz.v1.Authz/DeletePolicies", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authzClient) ListUserGroups(ctx context.Context, in *types.UserGroup, opts ...grpc.CallOption) (*types.UserGroups, error) {
	out := new(types.UserGroups)
	err := c.cc.Invoke(ctx, "/rafay.dev.rpc.authz.v1.Authz/ListUserGroups", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authzClient) CreateUserGroups(ctx context.Context, in *types.UserGroups, opts ...grpc.CallOption) (*types.BoolReply, error) {
	out := new(types.BoolReply)
	err := c.cc.Invoke(ctx, "/rafay.dev.rpc.authz.v1.Authz/CreateUserGroups", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authzClient) DeleteUserGroups(ctx context.Context, in *types.UserGroup, opts ...grpc.CallOption) (*types.BoolReply, error) {
	out := new(types.BoolReply)
	err := c.cc.Invoke(ctx, "/rafay.dev.rpc.authz.v1.Authz/DeleteUserGroups", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authzClient) ListRolePermissionMappings(ctx context.Context, in *types.FilteredRolePermissionMapping, opts ...grpc.CallOption) (*types.RolePermissionMappingList, error) {
	out := new(types.RolePermissionMappingList)
	err := c.cc.Invoke(ctx, "/rafay.dev.rpc.authz.v1.Authz/ListRolePermissionMappings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authzClient) CreateRolePermissionMappings(ctx context.Context, in *types.RolePermissionMappingList, opts ...grpc.CallOption) (*types.BoolReply, error) {
	out := new(types.BoolReply)
	err := c.cc.Invoke(ctx, "/rafay.dev.rpc.authz.v1.Authz/CreateRolePermissionMappings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authzClient) DeleteRolePermissionMappings(ctx context.Context, in *types.FilteredRolePermissionMapping, opts ...grpc.CallOption) (*types.BoolReply, error) {
	out := new(types.BoolReply)
	err := c.cc.Invoke(ctx, "/rafay.dev.rpc.authz.v1.Authz/DeleteRolePermissionMappings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthzServer is the server API for Authz service.
// All implementations should embed UnimplementedAuthzServer
// for forward compatibility
type AuthzServer interface {
	Enforce(context.Context, *types.EnforceRequest) (*types.BoolReply, error)
	// List Policies accpets Policy whose fileds are used for filtering
	// Filtering is done per field for the policy
	// For Example:
	// The Policy obj:
	//    sub => ""
	//    ns => ""
	//    proj => project1
	//    org => org1
	//    obj => ""
	//    act => ""
	// Returns policies related to project1 and org1 (Empty string matches all)
	ListPolicies(context.Context, *types.Policy) (*types.Policies, error)
	CreatePolicies(context.Context, *types.Policies) (*types.BoolReply, error)
	DeletePolicies(context.Context, *types.Policy) (*types.BoolReply, error)
	ListUserGroups(context.Context, *types.UserGroup) (*types.UserGroups, error)
	CreateUserGroups(context.Context, *types.UserGroups) (*types.BoolReply, error)
	DeleteUserGroups(context.Context, *types.UserGroup) (*types.BoolReply, error)
	ListRolePermissionMappings(context.Context, *types.FilteredRolePermissionMapping) (*types.RolePermissionMappingList, error)
	CreateRolePermissionMappings(context.Context, *types.RolePermissionMappingList) (*types.BoolReply, error)
	DeleteRolePermissionMappings(context.Context, *types.FilteredRolePermissionMapping) (*types.BoolReply, error)
}

// UnimplementedAuthzServer should be embedded to have forward compatible implementations.
type UnimplementedAuthzServer struct {
}

func (UnimplementedAuthzServer) Enforce(context.Context, *types.EnforceRequest) (*types.BoolReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Enforce not implemented")
}
func (UnimplementedAuthzServer) ListPolicies(context.Context, *types.Policy) (*types.Policies, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPolicies not implemented")
}
func (UnimplementedAuthzServer) CreatePolicies(context.Context, *types.Policies) (*types.BoolReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePolicies not implemented")
}
func (UnimplementedAuthzServer) DeletePolicies(context.Context, *types.Policy) (*types.BoolReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePolicies not implemented")
}
func (UnimplementedAuthzServer) ListUserGroups(context.Context, *types.UserGroup) (*types.UserGroups, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUserGroups not implemented")
}
func (UnimplementedAuthzServer) CreateUserGroups(context.Context, *types.UserGroups) (*types.BoolReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUserGroups not implemented")
}
func (UnimplementedAuthzServer) DeleteUserGroups(context.Context, *types.UserGroup) (*types.BoolReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserGroups not implemented")
}
func (UnimplementedAuthzServer) ListRolePermissionMappings(context.Context, *types.FilteredRolePermissionMapping) (*types.RolePermissionMappingList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRolePermissionMappings not implemented")
}
func (UnimplementedAuthzServer) CreateRolePermissionMappings(context.Context, *types.RolePermissionMappingList) (*types.BoolReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRolePermissionMappings not implemented")
}
func (UnimplementedAuthzServer) DeleteRolePermissionMappings(context.Context, *types.FilteredRolePermissionMapping) (*types.BoolReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRolePermissionMappings not implemented")
}

// UnsafeAuthzServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthzServer will
// result in compilation errors.
type UnsafeAuthzServer interface {
	mustEmbedUnimplementedAuthzServer()
}

func RegisterAuthzServer(s grpc.ServiceRegistrar, srv AuthzServer) {
	s.RegisterService(&Authz_ServiceDesc, srv)
}

func _Authz_Enforce_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.EnforceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthzServer).Enforce(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rafay.dev.rpc.authz.v1.Authz/Enforce",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthzServer).Enforce(ctx, req.(*types.EnforceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authz_ListPolicies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.Policy)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthzServer).ListPolicies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rafay.dev.rpc.authz.v1.Authz/ListPolicies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthzServer).ListPolicies(ctx, req.(*types.Policy))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authz_CreatePolicies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.Policies)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthzServer).CreatePolicies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rafay.dev.rpc.authz.v1.Authz/CreatePolicies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthzServer).CreatePolicies(ctx, req.(*types.Policies))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authz_DeletePolicies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.Policy)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthzServer).DeletePolicies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rafay.dev.rpc.authz.v1.Authz/DeletePolicies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthzServer).DeletePolicies(ctx, req.(*types.Policy))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authz_ListUserGroups_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.UserGroup)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthzServer).ListUserGroups(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rafay.dev.rpc.authz.v1.Authz/ListUserGroups",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthzServer).ListUserGroups(ctx, req.(*types.UserGroup))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authz_CreateUserGroups_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.UserGroups)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthzServer).CreateUserGroups(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rafay.dev.rpc.authz.v1.Authz/CreateUserGroups",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthzServer).CreateUserGroups(ctx, req.(*types.UserGroups))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authz_DeleteUserGroups_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.UserGroup)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthzServer).DeleteUserGroups(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rafay.dev.rpc.authz.v1.Authz/DeleteUserGroups",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthzServer).DeleteUserGroups(ctx, req.(*types.UserGroup))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authz_ListRolePermissionMappings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.FilteredRolePermissionMapping)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthzServer).ListRolePermissionMappings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rafay.dev.rpc.authz.v1.Authz/ListRolePermissionMappings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthzServer).ListRolePermissionMappings(ctx, req.(*types.FilteredRolePermissionMapping))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authz_CreateRolePermissionMappings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.RolePermissionMappingList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthzServer).CreateRolePermissionMappings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rafay.dev.rpc.authz.v1.Authz/CreateRolePermissionMappings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthzServer).CreateRolePermissionMappings(ctx, req.(*types.RolePermissionMappingList))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authz_DeleteRolePermissionMappings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.FilteredRolePermissionMapping)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthzServer).DeleteRolePermissionMappings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rafay.dev.rpc.authz.v1.Authz/DeleteRolePermissionMappings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthzServer).DeleteRolePermissionMappings(ctx, req.(*types.FilteredRolePermissionMapping))
	}
	return interceptor(ctx, in, info, handler)
}

// Authz_ServiceDesc is the grpc.ServiceDesc for Authz service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Authz_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rafay.dev.rpc.authz.v1.Authz",
	HandlerType: (*AuthzServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Enforce",
			Handler:    _Authz_Enforce_Handler,
		},
		{
			MethodName: "ListPolicies",
			Handler:    _Authz_ListPolicies_Handler,
		},
		{
			MethodName: "CreatePolicies",
			Handler:    _Authz_CreatePolicies_Handler,
		},
		{
			MethodName: "DeletePolicies",
			Handler:    _Authz_DeletePolicies_Handler,
		},
		{
			MethodName: "ListUserGroups",
			Handler:    _Authz_ListUserGroups_Handler,
		},
		{
			MethodName: "CreateUserGroups",
			Handler:    _Authz_CreateUserGroups_Handler,
		},
		{
			MethodName: "DeleteUserGroups",
			Handler:    _Authz_DeleteUserGroups_Handler,
		},
		{
			MethodName: "ListRolePermissionMappings",
			Handler:    _Authz_ListRolePermissionMappings_Handler,
		},
		{
			MethodName: "CreateRolePermissionMappings",
			Handler:    _Authz_CreateRolePermissionMappings_Handler,
		},
		{
			MethodName: "DeleteRolePermissionMappings",
			Handler:    _Authz_DeleteRolePermissionMappings_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/rpc/v1/authz.proto",
}
