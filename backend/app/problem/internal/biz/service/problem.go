package service

import (
	"context"
	"github.com/pengdahong1225/oj-server/backend/app/common/errs"
	"github.com/pengdahong1225/oj-server/backend/app/problem/internal/repository/domain"
	"github.com/pengdahong1225/oj-server/backend/app/problem/internal/repository/model"
	"github.com/pengdahong1225/oj-server/backend/module/utils"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProblemService struct {
	pb.UnimplementedProblemServiceServer
	db *domain.MysqlDB
}

func NewProblemService() *ProblemService {
	var err error
	s := &ProblemService{}
	s.db, err = domain.NewMysqlDB()
	if err != nil {
		logrus.Fatalf("NewProblemService failed, err:%s", err.Error())
	}
	return s
}

func (ps *ProblemService) CreateProblem(ctx context.Context, in *pb.CreateProblemRequest) (*pb.CreateProblemResponse, error) {
	resp := &pb.CreateProblemResponse{}

	problem := &model.Problem{
		Title:        in.Title,
		Level:        in.Level,
		Tags:         []byte(utils.SpliceStringWithX(in.Tags, "#")),
		Description:  in.Description,
		CreateBy:     0,
		CommentCount: 0,
		Config:       nil,
	}

	id, err := ps.db.CreateProblem(problem)
	if err != nil {
		logrus.Errorf("CreateProblem failed, err:%s", err.Error())
		return nil, errs.UpdateFailed
	}
	resp.Id = id
	return resp, nil
}

func (ps *ProblemService) UploadConfig(ctx context.Context, in *pb.UploadConfigRequest) (*pb.UploadConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadConfig not implemented")
}
func (ps *ProblemService) PublishProblem(ctx context.Context, in *pb.PublishProblemRequest) (*pb.PublishProblemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishProblem not implemented")
}

func (ps *ProblemService) UpdateProblem(ctx context.Context, in *pb.UpdateProblemRequest) (*pb.UpdateProblemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (ps *ProblemService) DeleteProblem(ctx context.Context, in *pb.DeleteProblemRequest) (*pb.DeleteProblemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (ps *ProblemService) GetProblemList(ctx context.Context, in *pb.GetProblemListRequest) (*pb.GetProblemListResponse, error) {
	resp := &pb.GetProblemListResponse{}

	total, problems, err := ps.db.QueryProblemList(int(in.Page), int(in.PageSize), in.Keyword, in.Tag)
	if err != nil {
		return nil, err
	}
	resp.Total = int32(total)
	resp.Data = model.TransformList(problems)
	return resp, nil
}

func (ps *ProblemService) GetProblemData(ctx context.Context, in *pb.GetProblemRequest) (*pb.GetProblemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (ps *ProblemService) GetTagList(ctx context.Context, empty *emptypb.Empty) (*pb.GetTagListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (ps *ProblemService) SubmitProblem(ctx context.Context, in *pb.SubmitProblemRequest) (*pb.SubmitProblemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (ps *ProblemService) GetSubmitRecordList(ctx context.Context, in *pb.GetSubmitRecordListRequest) (*pb.GetSubmitRecordListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (ps *ProblemService) GetSubmitRecordData(ctx context.Context, in *pb.GetSubmitRecordRequest) (*pb.GetSubmitRecordResponse, error) {
	//TODO implement me
	panic("implement me")
}
