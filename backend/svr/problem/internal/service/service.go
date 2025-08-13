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
	"oj-server/global"
	"oj-server/module/db/model"
	"oj-server/module/mq"
	"oj-server/proto/pb"
	"oj-server/svr/problem/internal/biz"
	"oj-server/svr/problem/internal/data"
	"oj-server/utils"
	"os"
	"path/filepath"
)

// 题目服务
type ProblemService struct {
	pb.UnimplementedProblemServiceServer
	pc *biz.ProblemUseCase
	rc *biz.RecordUseCase

	problem_producer *mq.Producer // 判题任务生产者
	comment_consumer *mq.Consumer // 评论任务消费者
}

func NewProblemService() *ProblemService {
	var err error
	s := &ProblemService{}

	pr, err := data.NewProblemRepo()
	if err != nil {
		logrus.Fatalf("NewProblemService failed, err:%s", err.Error())
	}
	rr, err := data.NewRecordRepo()
	if err != nil {
		logrus.Fatalf("NewProblemService failed, err:%s", err.Error())
	}

	s.pc = biz.NewProblemUseCase(pr) // 注入实现
	s.rc = biz.NewRecordUseCase(rr)  // 注入实现

	s.problem_producer = mq.NewProducer(
		global.RabbitMqExchangeKind,
		global.RabbitMqExchangeName,
		global.RabbitMqJudgeQueue,
		global.RabbitMqJudgeKey,
	)
	s.comment_consumer = mq.NewConsumer(
		global.RabbitMqExchangeKind,
		global.RabbitMqExchangeName,
		global.RabbitMqCommentQueue,
		global.RabbitMqCommentKey,
		"", // 消费者标签，用于区别不同的消费者
	)

	return s
}

func (ps *ProblemService) PublishSubmit2MQ(data []byte) bool {
	return ps.problem_producer.Publish(data)
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

	id, err := ps.pc.CreateProblem(problem)
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
			filePath := fmt.Sprintf("%s/%d.json", global.ProblemConfigPath, problemID)

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
		FilePath: fmt.Sprintf("%s/%d.json", global.ProblemConfigPath, problemID),
		Size:     fileSize,
	})
}
func (ps *ProblemService) PublishProblem(ctx context.Context, in *pb.PublishProblemRequest) (*pb.PublishProblemResponse, error) {
	resp := &pb.PublishProblemResponse{}

	err := ps.pc.UpdateProblemStatus(in.Id, 1)
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
	err := ps.pc.UpdateProblem(problem)
	if err != nil {
		logrus.Errorf("UpdateProblem failed, err:%s", err.Error())
		return nil, err
	}
	return resp, nil
}

func (ps *ProblemService) DeleteProblem(ctx context.Context, in *pb.DeleteProblemRequest) (*pb.DeleteProblemResponse, error) {
	resp := &pb.DeleteProblemResponse{}
	return resp, ps.pc.DeleteProblem(in.Id)
}

func (ps *ProblemService) GetProblemList(ctx context.Context, in *pb.GetProblemListRequest) (*pb.GetProblemListResponse, error) {
	resp := &pb.GetProblemListResponse{}

	total, problems, err := ps.pc.QueryProblemList(int(in.Page), int(in.PageSize), in.Keyword, in.Tag)
	if err != nil {
		return nil, err
	}
	resp.Total = int32(total)
	resp.Data = model.TransformList(problems)
	return resp, nil
}

func (ps *ProblemService) GetProblemData(ctx context.Context, in *pb.GetProblemRequest) (*pb.GetProblemResponse, error) {
	resp := &pb.GetProblemResponse{}

	problem, err := ps.pc.QueryProblemData(in.Id)
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
	key := fmt.Sprintf("%s:%d", global.UserLockPrefix, uid)
	_, err := ps.rc.Lock(key, global.UserLockTTL)
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
	form_data, err := proto.Marshal(&form)
	if err != nil {
		logrus.Errorf("marshal failed:%s", err.Error())
		_ = ps.rc.UnLock(key) // 释放锁
		return nil, status.Errorf(codes.Internal, "marshal failed")
	}
	// 提交到mq
	if !ps.PublishSubmit2MQ(form_data) {
		logrus.Errorf("publish to mq failed")
		_ = ps.rc.UnLock(key) // 释放锁
		return nil, status.Errorf(codes.Internal, "服务器错误")
	}

	resp.TaskId = taskId
	resp.Message = "题目提交成功"
	return resp, nil
}

func (ps *ProblemService) GetTagList(ctx context.Context, empty *emptypb.Empty) (*pb.GetTagListResponse, error) {
	resp := &pb.GetTagListResponse{}

	// 查询所有的标签
	tagList, err := ps.pc.QueryTagList()
	if err != nil {
		return nil, err
	}
	resp.Data = tagList
	return resp, nil
}
