package service

import (
	"context"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProblemService struct {
	pb.UnimplementedProblemServiceServer
}

func (ProblemService) UpdateProblem(ctx context.Context, request *pb.UpdateProblemRequest) (*pb.UpdateProblemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (ProblemService) DeleteProblem(ctx context.Context, request *pb.DeleteProblemRequest) (*pb.DeleteProblemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (ProblemService) GetProblemList(ctx context.Context, request *pb.GetProblemListRequest) (*pb.GetProblemListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (ProblemService) GetProblemData(ctx context.Context, request *pb.GetProblemRequest) (*pb.GetProblemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (ProblemService) GetTagList(ctx context.Context, empty *emptypb.Empty) (*pb.GetTagListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (ProblemService) SubmitProblem(ctx context.Context, request *pb.SubmitProblemRequest) (*pb.SubmitProblemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (ProblemService) GetSubmitRecordList(ctx context.Context, request *pb.GetSubmitRecordListRequest) (*pb.GetSubmitRecordListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (ProblemService) GetSubmitRecordData(ctx context.Context, request *pb.GetSubmitRecordRequest) (*pb.GetSubmitRecordResponse, error) {
	//TODO implement me
	panic("implement me")
}
