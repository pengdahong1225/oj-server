package handler

import (
	"context"
	pb2 "db-service/internal/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

func (receiver *DBServiceServer) GetUserSubmitRecord(ctx context.Context, request *pb2.GetUserSubmitRecordRequest) (*pb2.GetUserSubmitRecordResponse, error) {
	return nil, nil
}

func (receiver *DBServiceServer) UpdateUserSubmitRecord(ctx context.Context, request *pb2.UpdateUserSubmitRecordRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
