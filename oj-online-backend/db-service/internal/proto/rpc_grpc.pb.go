// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.10.0
// source: rpc.proto

package pb

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	DBService_GetUserDataByMobile_FullMethodName    = "/DBService/GetUserDataByMobile"
	DBService_GetUserDataByUid_FullMethodName       = "/DBService/GetUserDataByUid"
	DBService_CreateUserData_FullMethodName         = "/DBService/CreateUserData"
	DBService_UpdateUserData_FullMethodName         = "/DBService/UpdateUserData"
	DBService_DeleteUserData_FullMethodName         = "/DBService/DeleteUserData"
	DBService_GetUserList_FullMethodName            = "/DBService/GetUserList"
	DBService_GetUserSolvedList_FullMethodName      = "/DBService/GetUserSolvedList"
	DBService_GetProblemData_FullMethodName         = "/DBService/GetProblemData"
	DBService_CreateProblemData_FullMethodName      = "/DBService/CreateProblemData"
	DBService_UpdateProblemData_FullMethodName      = "/DBService/UpdateProblemData"
	DBService_DeleteProblemData_FullMethodName      = "/DBService/DeleteProblemData"
	DBService_GetProblemList_FullMethodName         = "/DBService/GetProblemList"
	DBService_QueryProblemWithName_FullMethodName   = "/DBService/QueryProblemWithName"
	DBService_GetUserSubmitRecord_FullMethodName    = "/DBService/GetUserSubmitRecord"
	DBService_UpdateUserSubmitRecord_FullMethodName = "/DBService/UpdateUserSubmitRecord"
)

// DBServiceClient is the client API for DBService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DBServiceClient interface {
	// user接口
	GetUserDataByMobile(ctx context.Context, in *GetUserDataByMobileRequest, opts ...grpc.CallOption) (*GetUserResponse, error)
	GetUserDataByUid(ctx context.Context, in *GetUserDataByUidRequest, opts ...grpc.CallOption) (*GetUserResponse, error)
	CreateUserData(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	UpdateUserData(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	DeleteUserData(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	GetUserList(ctx context.Context, in *GetUserListRequest, opts ...grpc.CallOption) (*GetUserListResponse, error)
	GetUserSolvedList(ctx context.Context, in *GetUserSolvedListRequest, opts ...grpc.CallOption) (*GetUserSolvedListResponse, error)
	// problem接口
	GetProblemData(ctx context.Context, in *GetProblemRequest, opts ...grpc.CallOption) (*GetProblemResponse, error)
	CreateProblemData(ctx context.Context, in *CreateProblemRequest, opts ...grpc.CallOption) (*CreateProblemResponse, error)
	UpdateProblemData(ctx context.Context, in *UpdateProblemRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	DeleteProblemData(ctx context.Context, in *DeleteProblemRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	GetProblemList(ctx context.Context, in *GetProblemListRequest, opts ...grpc.CallOption) (*GetProblemListResponse, error)
	QueryProblemWithName(ctx context.Context, in *QueryProblemWithNameRequest, opts ...grpc.CallOption) (*QueryProblemWithNameResponse, error)
	GetUserSubmitRecord(ctx context.Context, in *GetUserSubmitRecordRequest, opts ...grpc.CallOption) (*GetUserSubmitRecordResponse, error)
	UpdateUserSubmitRecord(ctx context.Context, in *UpdateUserSubmitRecordRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type dBServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDBServiceClient(cc grpc.ClientConnInterface) DBServiceClient {
	return &dBServiceClient{cc}
}

func (c *dBServiceClient) GetUserDataByMobile(ctx context.Context, in *GetUserDataByMobileRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	out := new(GetUserResponse)
	err := c.cc.Invoke(ctx, DBService_GetUserDataByMobile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) GetUserDataByUid(ctx context.Context, in *GetUserDataByUidRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	out := new(GetUserResponse)
	err := c.cc.Invoke(ctx, DBService_GetUserDataByUid_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) CreateUserData(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, DBService_CreateUserData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) UpdateUserData(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, DBService_UpdateUserData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) DeleteUserData(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, DBService_DeleteUserData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) GetUserList(ctx context.Context, in *GetUserListRequest, opts ...grpc.CallOption) (*GetUserListResponse, error) {
	out := new(GetUserListResponse)
	err := c.cc.Invoke(ctx, DBService_GetUserList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) GetUserSolvedList(ctx context.Context, in *GetUserSolvedListRequest, opts ...grpc.CallOption) (*GetUserSolvedListResponse, error) {
	out := new(GetUserSolvedListResponse)
	err := c.cc.Invoke(ctx, DBService_GetUserSolvedList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) GetProblemData(ctx context.Context, in *GetProblemRequest, opts ...grpc.CallOption) (*GetProblemResponse, error) {
	out := new(GetProblemResponse)
	err := c.cc.Invoke(ctx, DBService_GetProblemData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) CreateProblemData(ctx context.Context, in *CreateProblemRequest, opts ...grpc.CallOption) (*CreateProblemResponse, error) {
	out := new(CreateProblemResponse)
	err := c.cc.Invoke(ctx, DBService_CreateProblemData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) UpdateProblemData(ctx context.Context, in *UpdateProblemRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, DBService_UpdateProblemData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) DeleteProblemData(ctx context.Context, in *DeleteProblemRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, DBService_DeleteProblemData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) GetProblemList(ctx context.Context, in *GetProblemListRequest, opts ...grpc.CallOption) (*GetProblemListResponse, error) {
	out := new(GetProblemListResponse)
	err := c.cc.Invoke(ctx, DBService_GetProblemList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) QueryProblemWithName(ctx context.Context, in *QueryProblemWithNameRequest, opts ...grpc.CallOption) (*QueryProblemWithNameResponse, error) {
	out := new(QueryProblemWithNameResponse)
	err := c.cc.Invoke(ctx, DBService_QueryProblemWithName_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) GetUserSubmitRecord(ctx context.Context, in *GetUserSubmitRecordRequest, opts ...grpc.CallOption) (*GetUserSubmitRecordResponse, error) {
	out := new(GetUserSubmitRecordResponse)
	err := c.cc.Invoke(ctx, DBService_GetUserSubmitRecord_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) UpdateUserSubmitRecord(ctx context.Context, in *UpdateUserSubmitRecordRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, DBService_UpdateUserSubmitRecord_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DBServiceServer is the server API for DBService service.
// All implementations must embed UnimplementedDBServiceServer
// for forward compatibility
type DBServiceServer interface {
	// user接口
	GetUserDataByMobile(context.Context, *GetUserDataByMobileRequest) (*GetUserResponse, error)
	GetUserDataByUid(context.Context, *GetUserDataByUidRequest) (*GetUserResponse, error)
	CreateUserData(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	UpdateUserData(context.Context, *UpdateUserRequest) (*empty.Empty, error)
	DeleteUserData(context.Context, *DeleteUserRequest) (*empty.Empty, error)
	GetUserList(context.Context, *GetUserListRequest) (*GetUserListResponse, error)
	GetUserSolvedList(context.Context, *GetUserSolvedListRequest) (*GetUserSolvedListResponse, error)
	// problem接口
	GetProblemData(context.Context, *GetProblemRequest) (*GetProblemResponse, error)
	CreateProblemData(context.Context, *CreateProblemRequest) (*CreateProblemResponse, error)
	UpdateProblemData(context.Context, *UpdateProblemRequest) (*empty.Empty, error)
	DeleteProblemData(context.Context, *DeleteProblemRequest) (*empty.Empty, error)
	GetProblemList(context.Context, *GetProblemListRequest) (*GetProblemListResponse, error)
	QueryProblemWithName(context.Context, *QueryProblemWithNameRequest) (*QueryProblemWithNameResponse, error)
	GetUserSubmitRecord(context.Context, *GetUserSubmitRecordRequest) (*GetUserSubmitRecordResponse, error)
	UpdateUserSubmitRecord(context.Context, *UpdateUserSubmitRecordRequest) (*empty.Empty, error)
	mustEmbedUnimplementedDBServiceServer()
}

// UnimplementedDBServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDBServiceServer struct {
}

func (UnimplementedDBServiceServer) GetUserDataByMobile(context.Context, *GetUserDataByMobileRequest) (*GetUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserDataByMobile not implemented")
}
func (UnimplementedDBServiceServer) GetUserDataByUid(context.Context, *GetUserDataByUidRequest) (*GetUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserDataByUid not implemented")
}
func (UnimplementedDBServiceServer) CreateUserData(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUserData not implemented")
}
func (UnimplementedDBServiceServer) UpdateUserData(context.Context, *UpdateUserRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserData not implemented")
}
func (UnimplementedDBServiceServer) DeleteUserData(context.Context, *DeleteUserRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserData not implemented")
}
func (UnimplementedDBServiceServer) GetUserList(context.Context, *GetUserListRequest) (*GetUserListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserList not implemented")
}
func (UnimplementedDBServiceServer) GetUserSolvedList(context.Context, *GetUserSolvedListRequest) (*GetUserSolvedListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserSolvedList not implemented")
}
func (UnimplementedDBServiceServer) GetProblemData(context.Context, *GetProblemRequest) (*GetProblemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProblemData not implemented")
}
func (UnimplementedDBServiceServer) CreateProblemData(context.Context, *CreateProblemRequest) (*CreateProblemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProblemData not implemented")
}
func (UnimplementedDBServiceServer) UpdateProblemData(context.Context, *UpdateProblemRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProblemData not implemented")
}
func (UnimplementedDBServiceServer) DeleteProblemData(context.Context, *DeleteProblemRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteProblemData not implemented")
}
func (UnimplementedDBServiceServer) GetProblemList(context.Context, *GetProblemListRequest) (*GetProblemListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProblemList not implemented")
}
func (UnimplementedDBServiceServer) QueryProblemWithName(context.Context, *QueryProblemWithNameRequest) (*QueryProblemWithNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryProblemWithName not implemented")
}
func (UnimplementedDBServiceServer) GetUserSubmitRecord(context.Context, *GetUserSubmitRecordRequest) (*GetUserSubmitRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserSubmitRecord not implemented")
}
func (UnimplementedDBServiceServer) UpdateUserSubmitRecord(context.Context, *UpdateUserSubmitRecordRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserSubmitRecord not implemented")
}
func (UnimplementedDBServiceServer) mustEmbedUnimplementedDBServiceServer() {}

// UnsafeDBServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DBServiceServer will
// result in compilation errors.
type UnsafeDBServiceServer interface {
	mustEmbedUnimplementedDBServiceServer()
}

func RegisterDBServiceServer(s grpc.ServiceRegistrar, srv DBServiceServer) {
	s.RegisterService(&DBService_ServiceDesc, srv)
}

func _DBService_GetUserDataByMobile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserDataByMobileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).GetUserDataByMobile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_GetUserDataByMobile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).GetUserDataByMobile(ctx, req.(*GetUserDataByMobileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_GetUserDataByUid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserDataByUidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).GetUserDataByUid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_GetUserDataByUid_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).GetUserDataByUid(ctx, req.(*GetUserDataByUidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_CreateUserData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).CreateUserData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_CreateUserData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).CreateUserData(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_UpdateUserData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).UpdateUserData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_UpdateUserData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).UpdateUserData(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_DeleteUserData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).DeleteUserData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_DeleteUserData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).DeleteUserData(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_GetUserList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).GetUserList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_GetUserList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).GetUserList(ctx, req.(*GetUserListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_GetUserSolvedList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserSolvedListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).GetUserSolvedList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_GetUserSolvedList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).GetUserSolvedList(ctx, req.(*GetUserSolvedListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_GetProblemData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProblemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).GetProblemData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_GetProblemData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).GetProblemData(ctx, req.(*GetProblemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_CreateProblemData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateProblemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).CreateProblemData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_CreateProblemData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).CreateProblemData(ctx, req.(*CreateProblemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_UpdateProblemData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProblemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).UpdateProblemData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_UpdateProblemData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).UpdateProblemData(ctx, req.(*UpdateProblemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_DeleteProblemData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteProblemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).DeleteProblemData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_DeleteProblemData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).DeleteProblemData(ctx, req.(*DeleteProblemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_GetProblemList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProblemListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).GetProblemList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_GetProblemList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).GetProblemList(ctx, req.(*GetProblemListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_QueryProblemWithName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryProblemWithNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).QueryProblemWithName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_QueryProblemWithName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).QueryProblemWithName(ctx, req.(*QueryProblemWithNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_GetUserSubmitRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserSubmitRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).GetUserSubmitRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_GetUserSubmitRecord_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).GetUserSubmitRecord(ctx, req.(*GetUserSubmitRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_UpdateUserSubmitRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserSubmitRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).UpdateUserSubmitRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_UpdateUserSubmitRecord_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).UpdateUserSubmitRecord(ctx, req.(*UpdateUserSubmitRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DBService_ServiceDesc is the grpc.ServiceDesc for DBService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DBService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "DBService",
	HandlerType: (*DBServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserDataByMobile",
			Handler:    _DBService_GetUserDataByMobile_Handler,
		},
		{
			MethodName: "GetUserDataByUid",
			Handler:    _DBService_GetUserDataByUid_Handler,
		},
		{
			MethodName: "CreateUserData",
			Handler:    _DBService_CreateUserData_Handler,
		},
		{
			MethodName: "UpdateUserData",
			Handler:    _DBService_UpdateUserData_Handler,
		},
		{
			MethodName: "DeleteUserData",
			Handler:    _DBService_DeleteUserData_Handler,
		},
		{
			MethodName: "GetUserList",
			Handler:    _DBService_GetUserList_Handler,
		},
		{
			MethodName: "GetUserSolvedList",
			Handler:    _DBService_GetUserSolvedList_Handler,
		},
		{
			MethodName: "GetProblemData",
			Handler:    _DBService_GetProblemData_Handler,
		},
		{
			MethodName: "CreateProblemData",
			Handler:    _DBService_CreateProblemData_Handler,
		},
		{
			MethodName: "UpdateProblemData",
			Handler:    _DBService_UpdateProblemData_Handler,
		},
		{
			MethodName: "DeleteProblemData",
			Handler:    _DBService_DeleteProblemData_Handler,
		},
		{
			MethodName: "GetProblemList",
			Handler:    _DBService_GetProblemList_Handler,
		},
		{
			MethodName: "QueryProblemWithName",
			Handler:    _DBService_QueryProblemWithName_Handler,
		},
		{
			MethodName: "GetUserSubmitRecord",
			Handler:    _DBService_GetUserSubmitRecord_Handler,
		},
		{
			MethodName: "UpdateUserSubmitRecord",
			Handler:    _DBService_UpdateUserSubmitRecord_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc.proto",
}
