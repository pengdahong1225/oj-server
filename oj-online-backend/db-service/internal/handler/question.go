package handler

import (
	"context"
	pb "db-service/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (receiver *DBServiceServer) GetQuestionData(ctx context.Context, request *pb.GetQuestionRequest) (*pb.GetQuestionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetQuestionData not implemented")
}
func (receiver *DBServiceServer) CreateQuestionData(ctx context.Context, request *pb.CreateQuestionRequest) (*pb.CreateQuestionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateQuestionData not implemented")
}
func (receiver *DBServiceServer) UpdateQuestionData(ctx context.Context, request *pb.UpdateQuestionRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateQuestionData not implemented")
}
func (receiver *DBServiceServer) DeleteQuestionData(ctx context.Context, request *pb.DeleteQuestionRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteQuestionData not implemented")
}
func (receiver *DBServiceServer) GetQuestionList(ctx context.Context, request *pb.GetQuestionListRequest) (*pb.GetQuestionListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetQuestionList not implemented")
}
