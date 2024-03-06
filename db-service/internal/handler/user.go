package handler

import (
	"context"
	pb "db-service/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DBServiceServer struct {
	pb.UnimplementedDBServiceServer
}

func (receiver *DBServiceServer) GetUserData(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserData not implemented")
}

func (receiver *DBServiceServer) CreateUserData(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUserData not implemented")
}
func (receiver *DBServiceServer) UpdateUserData(ctx context.Context, request *pb.UpdateUserRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserData not implemented")
}

func (receiver *DBServiceServer) DeleteUserData(ctx context.Context, request *pb.DeleteUserRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserData not implemented")
}
func (receiver *DBServiceServer) GetUserList(ctx context.Context, request *pb.GetUserListRequest) (*pb.GetUserListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserList not implemented")
}
