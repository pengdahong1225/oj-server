package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"oj-server/global"
	"oj-server/pkg/mq"
	"oj-server/pkg/proto/pb"
	"oj-server/pkg/registry"
	"oj-server/svr/problem/internal/biz"
	"oj-server/svr/problem/internal/configs"
	"oj-server/svr/problem/internal/data"
	"oj-server/svr/problem/internal/model"
	"os"
	"path/filepath"
	"strconv"
)

// 题目服务
type ProblemService struct {
	pb.UnimplementedProblemServiceServer
	uc *biz.ProblemUseCase

	problem_producer *mq.Producer // 判题任务生产者
}

func NewProblemService() *ProblemService {
	repo, err := data.NewProblemRepo()
	if err != nil {
		logrus.Fatalf("NewProblemService failed, err:%s", err.Error())
	}

	mqCfg := configs.AppConf.MQCfg
	amqpClient := mq.NewClient(
		&mq.Options{
			Host:     mqCfg.Host,
			Port:     mqCfg.Port,
			User:     mqCfg.User,
			PassWord: mqCfg.PassWord,
			VHost:    mqCfg.VHost,
		},
	)

	return &ProblemService{
		uc: biz.NewProblemUseCase(repo), // 注入实现
		problem_producer: &mq.Producer{
			AmqpClient: amqpClient, // 注入client
			ExKind:     global.RabbitMqExchangeKind,
			ExName:     global.RabbitMqExchangeName,

			QueName:    global.RabbitMqJudgeSubmitQueue,
			RoutingKey: global.RabbitMqJudgeSubmitKey,
		},
	}
}

func (ps *ProblemService) CreateProblem(ctx context.Context, in *pb.CreateProblemRequest) (*pb.CreateProblemResponse, error) {
	resp := &pb.CreateProblemResponse{}

	problem := &model.Problem{
		Title:       in.Title,
		Level:       int8(in.Level),
		Description: in.Description,
	}
	var err error
	problem.Tags, err = json.Marshal(in.Tags)
	if err != nil {
		logrus.Errorf("json marshal failed:%s", err.Error())
		return nil, err
	}

	id, err := ps.uc.CreateProblem(problem)
	if err != nil {
		logrus.Errorf("CreateProblem failed, err:%s", err.Error())
		return nil, err
	}
	resp.Id = id
	return resp, nil
}
func (ps *ProblemService) UploadConfig(stream pb.ProblemService_UploadConfigServer) error {
	var (
		saveSize int64
		filePath string
		writer   io.Writer
	)
	// 获取元数据
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Error(codes.Unauthenticated, "no metadata found")
	}
	mds := md.Get("problem_id")
	if len(mds) == 0 {
		return status.Error(codes.Unauthenticated, "problem_id missing")
	}
	problemId, _ := strconv.Atoi(mds[0]) // 转回 int

	// 查询题目是否存在
	problem, err := ps.uc.QueryProblemData(int64(problemId))
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return err
	}

	// 创建目录
	filePath = fmt.Sprintf("%s/%d.json", global.ProblemConfigPath, problemId)
	if err = os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
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

		// 写入分片
		if n, err := writer.Write(chunk.Content); err != nil {
			logrus.Errorf("write chunk failed: %v", err)
			return status.Error(codes.Internal, "write chunk failed")
		} else {
			saveSize += int64(n)
		}
	}

	// 更新数据库
	if err = ps.uc.UpdateProblemConfig(problem.ID, filePath); err != nil {
		logrus.Errorf("UpdateProblem failed, err:%s", err.Error())
		return status.Error(codes.Internal, "update problem failed")
	}

	// 返回成功响应
	return stream.SendAndClose(&pb.UploadConfigResponse{
		FilePath: filePath,
		Size:     saveSize,
	})
}
func (ps *ProblemService) PublishProblem(ctx context.Context, in *pb.PublishProblemRequest) (*emptypb.Empty, error) {
	if err := ps.uc.UpdateProblemStatus(in.Id, 1); err != nil {
		return nil, err
	}
	return nil, nil
}
func (ps *ProblemService) HideProblem(ctx context.Context, in *pb.HideProblemRequest) (*emptypb.Empty, error) {
	if err := ps.uc.UpdateProblemStatus(in.Id, 0); err != nil {
		return nil, err
	}
	return nil, nil
}

// 更新题目基础信息
// 标题、等级、标签、描述、创建者、状态
func (ps *ProblemService) UpdateProblem(ctx context.Context, in *pb.UpdateProblemRequest) (*emptypb.Empty, error) {
	problem := &model.Problem{
		ID:          in.Data.Id,
		Title:       in.Data.Title,
		Level:       int8(in.Data.Level),
		Description: in.Data.Description,
	}
	var err error
	problem.Tags, err = json.Marshal(in.Data.Tags)
	if err != nil {
		logrus.Errorf("json marshal failed:%s", err.Error())
		return nil, err
	}

	if err = ps.uc.UpdateProblem(problem); err != nil {
		logrus.Errorf("UpdateProblem failed, err:%s", err.Error())
		return nil, err
	}
	return nil, nil
}

func (ps *ProblemService) DeleteProblem(ctx context.Context, in *pb.DeleteProblemRequest) (*emptypb.Empty, error) {
	return nil, ps.uc.DeleteProblem(in.Id)
}

func (ps *ProblemService) GetProblemList(ctx context.Context, in *pb.GetProblemListRequest) (*pb.GetProblemListResponse, error) {
	resp := &pb.GetProblemListResponse{}

	// 获取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no metadata found")
	}
	roles := md.Get("role")
	if len(roles) == 0 {
		return nil, status.Error(codes.Unauthenticated, "role missing")
	}
	role, _ := strconv.Atoi(roles[0]) // 转回 int

	total, problems, err := ps.uc.QueryProblemList(int(in.Page), int(in.PageSize), in.Keyword, in.Tag, int32(role))
	if err != nil {
		return nil, err
	}
	resp.Total = total
	for _, problem := range problems {
		pbProblem := &pb.Problem{
			Id:          problem.ID,
			CreateAt:    problem.CreateAt.Unix(),
			UpdateAt:    problem.UpdateAt.Unix(),
			Title:       problem.Title,
			Description: problem.Description,
			Level:       int32(problem.Level),
			Status:      int32(problem.Status),
		}
		if err = json.Unmarshal(problem.Tags, &pbProblem.Tags); err != nil {
			logrus.Errorf("json unmarshal failed:%s", err.Error())
		}
		resp.Data = append(resp.Data, pbProblem)
	}

	return resp, nil
}

func (ps *ProblemService) GetProblemDetail(ctx context.Context, in *pb.GetProblemRequest) (*pb.GetProblemResponse, error) {
	resp := &pb.GetProblemResponse{}

	problem, err := ps.uc.QueryProblemData(in.Id)
	if err != nil {
		return nil, err
	}
	resp.Problem = &pb.Problem{
		Id:          problem.ID,
		CreateAt:    problem.CreateAt.Unix(),
		UpdateAt:    problem.UpdateAt.Unix(),
		Title:       problem.Title,
		Description: problem.Description,
		Level:       int32(problem.Level),
		Status:      int32(problem.Status),
	}
	if err = json.Unmarshal(problem.Tags, &resp.Problem.Tags); err != nil {
		logrus.Errorf("json unmarshal failed:%s", err.Error())
	}

	return resp, nil
}

func (ps *ProblemService) SubmitProblem(ctx context.Context, in *pb.SubmitProblemRequest) (*pb.SubmitProblemResponse, error) {
	// 获取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no metadata found")
	}
	vals := md.Get("uid")
	if len(vals) == 0 {
		return nil, status.Error(codes.Unauthenticated, "uid missing")
	}
	uid, err := strconv.ParseInt(vals[0], 10, 64)
	if err != nil {
		logrus.Errorf("parse uid failed:%s", err.Error())
		return nil, err
	}

	taskId := fmt.Sprintf("%d_%d", uid, in.ProblemId)

	// 查询用户名
	username, err := ps.queryUserName(uid)
	if err != nil {
		logrus.Errorf("query user name failed:%s", err.Error())
		return nil, err
	}

	// 查询题目信息
	problem, err := ps.uc.QueryProblemData(in.ProblemId)
	if err != nil {
		logrus.Errorf("query problem config url failed:%s", err.Error())
		return nil, err
	}

	// 获取分布式锁
	key := fmt.Sprintf("%s:%d", global.UserLockPrefix, uid)
	_, err = ps.uc.Lock(key, global.UserLockTTL)
	if err != nil {
		logrus.Errorf("lock user failed:%s", err.Error())
		return nil, fmt.Errorf("判题中")
	}

	task := &pb.JudgeSubmission{
		Uid:         int64(uid),
		ProblemId:   in.ProblemId,
		Title:       in.Title,
		Lang:        in.Lang,
		Code:        in.Code,
		ConfigUrl:   problem.ConfigURL,
		TaskId:      taskId,
		Level:       int32(problem.Level),
		ProblemName: problem.Title,
		UserName:    username,
	}
	task_data, err := proto.Marshal(task)
	if err != nil {
		logrus.Errorf("marshal failed:%s", err.Error())
		_ = ps.uc.UnLock(key) // 释放锁
		return nil, status.Errorf(codes.Internal, "marshal failed")
	}
	// 提交到mq
	if err = ps.problem_producer.Publish(task_data); err != nil {
		logrus.Errorf("publish to mq failed: %v", err)
		_ = ps.uc.UnLock(key) // 释放锁
		return nil, status.Errorf(codes.Internal, "服务器错误")
	}

	return &pb.SubmitProblemResponse{
		TaskId: taskId,
	}, nil
}

func (ps *ProblemService) queryUserName(uid int64) (string, error) {
	// 调用用户服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.UserService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		return "", err
	}
	client := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), global.GrpcTimeout)
	defer cancel()
	resp, err := client.GetUserInfoByUid(ctx, &pb.GetUserInfoRequest{
		Uid: uid,
	})
	if err != nil {
		return "", err
	}
	return resp.Data.Nickname, nil
}

func (ps *ProblemService) GetTagList(ctx context.Context, empty *emptypb.Empty) (*pb.GetTagListResponse, error) {
	resp := &pb.GetTagListResponse{}

	// 查询所有的标签
	tagList, err := ps.uc.QueryTagList()
	if err != nil {
		return nil, err
	}
	resp.Data = tagList
	return resp, nil
}
