package service

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"oj-server/app/problem/internal/productor"
	"oj-server/app/problem/internal/repository/cache"
	"oj-server/app/problem/internal/repository/domain"
	"oj-server/consts"
	"oj-server/module/utils"
	"oj-server/proto/pb"
	"os"
	"path/filepath"
	"oj-server/module/model"
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
		return nil, err
	}
	resp.Id = id
	return resp, nil
}
func (ps *ProblemService) UploadConfig(stream pb.ProblemService_UploadConfigServer) error {
	var (
		problemID int64
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
			filePath := fmt.Sprintf("%s/%d.json", consts.ProblemConfigPath, problemID)

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
		FilePath: fmt.Sprintf("%s/%d.json", consts.ProblemConfigPath, problemID),
		Size:     fileSize,
	})
}
func (ps *ProblemService) PublishProblem(ctx context.Context, in *pb.PublishProblemRequest) (*pb.PublishProblemResponse, error) {
	resp := &pb.PublishProblemResponse{}

	err := ps.db.UpdateProblemStatus(in.Id, 1)
	if err != nil {
		return nil, err
	}
	resp.Id = in.Id
	resp.Result = "publish problem success"
	return resp, nil
}

// 更新题目基础信息
// 标题、等级、标签、描述、创建者、状态
func (ps *ProblemService) UpdateProblem(ctx context.Context, in *pb.UpdateProblemRequest) (*pb.UpdateProblemResponse, error) {
	resp := &pb.UpdateProblemResponse{}

	problem := &model.Problem{
		ID:          in.Data.Id,
		Title:       in.Data.Title,
		Level:       in.Data.Level,
		Tags:        []byte(utils.SpliceStringWithX(in.Data.Tags, "#")),
		Description: in.Data.Description,
		CreateBy:    in.Data.CreateBy,
		Status:      in.Data.Status,
	}
	err := ps.db.UpdateProblem(problem)
	if err != nil {
		logrus.Errorf("UpdateProblem failed, err:%s", err.Error())
		return nil, err
	}
	return resp, nil
}

func (ps *ProblemService) DeleteProblem(ctx context.Context, in *pb.DeleteProblemRequest) (*pb.DeleteProblemResponse, error) {
	resp := &pb.DeleteProblemResponse{}
	return resp, ps.db.DeleteProblem(in.Id)
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

func (ps *ProblemService) SubmitProblem(ctx context.Context, in *pb.SubmitProblemRequest) (*pb.SubmitProblemResponse, error) {
	// 获取元数据
	uid := ctx.Value("uid").(int64)
	taskId := fmt.Sprintf("%d_%d", uid, in.ProblemId)

	resp := &pb.SubmitProblemResponse{}

	// todo 加锁，在超时or判题服务处理完后 才释放锁
	// 重复提交时会触发加锁失败
	_, err := cache.LockUser(uid)
	if err != nil {
		logrus.Errorf("lock user failed:%s", err.Error())
		return nil, fmt.Errorf("判题中")
	}

	form := pb.SubmitForm{
		Uid:       uid,
		ProblemId: in.ProblemId,
		Title:     in.Title,
		Lang:      in.Lang,
		Code:      in.Code,
	}
	data, err := proto.Marshal(&form)
	if err != nil {
		logrus.Errorf("marshal failed:%s", err.Error())
		_ = cache.UnLockUser(uid) // 释放锁
		return nil, status.Errorf(codes.Internal, "marshal failed")
	}
	// 提交到mq
	if !productor.Publish(data) {
		logrus.Errorf("publish to mq failed")
		_ = cache.UnLockUser(uid) // 释放锁
		return nil, status.Errorf(codes.Internal, "服务器错误")
	}

	resp.TaskId = taskId
	resp.Message = "题目提交成功"
	return resp, nil
}

func (ps *ProblemService) GetTagList(ctx context.Context, empty *emptypb.Empty) (*pb.GetTagListResponse, error) {
	//TODO implement me
	panic("implement me")
}
