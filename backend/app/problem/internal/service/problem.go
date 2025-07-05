package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"oj-server/app/common/errs"
	"oj-server/app/problem/internal/repository/domain"
	"oj-server/app/problem/internal/repository/model"
	"oj-server/module/utils"
	"oj-server/proto/pb"
	"os"
	"path/filepath"
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
	}

	id, err := ps.db.CreateProblem(problem)
	if err != nil {
		logrus.Errorf("CreateProblem failed, err:%s", err.Error())
		return nil, errs.UpdateFailed
	}
	resp.Id = id
	return resp, nil
}
func (ps *ProblemService) UploadConfig(stream pb.ProblemService_UploadConfigServer) error {
	var (
		problemID int64
		filename  string
		fileSize  int64
		writer    io.Writer
	)
	// 创建目标文件
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Errorf("receive chunk failed: %v", err)
			return status.Error(codes.Internal, "receive chunk failed")
		}

		// 首次接收时初始化
		if writer == nil {
			problemID = chunk.ProblemId
			filename = chunk.FileName
			filePath := fmt.Sprintf("/data/problems/%d/%s", problemID, filename)

			// 创建目录
			if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
				logrus.Errorf("create dir failed: %v", err)
				return status.Error(codes.Internal, "create dir failed")
			}

			// 创建文件
			f, err := os.Create(filePath)
			if err != nil {
				logrus.Errorf("create file failed: %v", err)
				return status.Error(codes.Internal, "create file failed")
			}
			defer f.Close()
			writer = f
		}

		// 写入分片
		if n, err := writer.Write(chunk.Content); err != nil {
			logrus.Errorf("write chunk failed: %v", err)
			return status.Error(codes.Internal, "write chunk failed")
		} else {
			fileSize += int64(n)
		}
	}

	// 返回成功响应
	return stream.SendAndClose(&pb.UploadConfigResponse{
		FilePath: fmt.Sprintf("/data/problems/%d/%s", problemID, filename),
		Size:     fileSize,
	})
}
func (ps *ProblemService) PublishProblem(ctx context.Context, in *pb.PublishProblemRequest) (*pb.PublishProblemResponse, error) {
	err := ps.db.UpdateProblemStatus(in.Id, 1)
	if errors.As(err, &errs.NotFound) {
		return nil, status.Error(codes.NotFound, "problem not found")
	}
	return nil, status.Error(codes.Internal, "update problem status failed")
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
	resp := &pb.GetProblemResponse{}

	problem, err := ps.db.QueryProblemData(in.Id)
	if err != nil {
		return nil, err
	}
	resp.Data = problem.Transform()
	return resp, nil
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
