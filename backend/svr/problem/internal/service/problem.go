package service

import (
	"context"
	"encoding/json"
	"fmt"
	"oj-server/common/consts"
	"oj-server/common/proto/pb"
	"oj-server/pkg/mq"
	"oj-server/svr/problem/internal/client"
	"oj-server/svr/problem/internal/config"
	"oj-server/svr/problem/internal/data"
	"oj-server/svr/problem/internal/model"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// 题目服务
type ProblemService struct {
	pb.UnimplementedProblemServiceServer
	repo          *data.ProblemRepo
	judgeProducer *mq.Producer
	clients       *client.Clients
}

func NewProblemService(clients *client.Clients) (*ProblemService, error) {
	repo, err := data.NewProblemRepo()
	if err != nil {
		logrus.Errorf("NewProblemService failed, err:%s", err.Error())
		return nil, err
	}
	cfg := config.MyConfig.MQ
	amqpClient := mq.NewClient(
		&mq.Options{
			Host:     cfg.Host,
			Port:     cfg.Port,
			User:     cfg.User,
			PassWord: cfg.PassWord,
			VHost:    cfg.VHost,
		},
	)

	return &ProblemService{
		repo: repo, // 注入实现
		judgeProducer: &mq.Producer{
			AmqpClient: amqpClient, // 注入client
			ExKind:     consts.RabbitMqExchangeKind,
			ExName:     consts.RabbitMqExchangeName,

			QueName:    consts.RabbitMqJudgeSubmitQueue,
			RoutingKey: consts.RabbitMqJudgeSubmitKey,
		},
		clients: clients,
	}, nil
}

func (ps *ProblemService) CreateProblem(ctx context.Context, in *pb.CreateProblemRequest) (*pb.CreateProblemResponse, error) {
	resp := &pb.CreateProblemResponse{}

	problem := &model.Problem{
		Title:  in.Title,
		Level:  int8(in.Level),
		Status: model.ProblemStatusPending,
	}
	var err error
	problem.Tags, err = json.Marshal(in.Tags)
	if err != nil {
		logrus.Errorf("json marshal failed:%s", err.Error())
		return nil, err
	}

	err = ps.repo.Transaction(func(tx *gorm.DB) error {
		// 1. 入库拿 id
		if err := tx.Create(problem).Error; err != nil {
			return err
		}

		// 2. 创建题目目录
		dir := filepath.Join(consts.PackageBaseDir, problem.Title)
		if err := os.MkdirAll(filepath.Join(dir, "testdata"), 0755); err != nil {
			return err
		}

		// 3. 回填 package_dir
		return tx.Model(&model.Problem{}).Where("id=?", problem.ID).Update("package_dir", dir).Error
	})
	if err != nil {
		logrus.Errorf("CreateProblem failed, err:%s", err.Error())
		return nil, err
	}
	resp.Id = problem.ID

	return resp, nil
}

func (ps *ProblemService) PublishProblem(ctx context.Context, in *pb.PublishProblemRequest) (*emptypb.Empty, error) {
	if err := ps.repo.UpdateProblemStatus(in.Id, 1); err != nil {
		return nil, err
	}
	return nil, nil
}
func (ps *ProblemService) HideProblem(ctx context.Context, in *pb.HideProblemRequest) (*emptypb.Empty, error) {
	if err := ps.repo.UpdateProblemStatus(in.Id, 0); err != nil {
		return nil, err
	}
	return nil, nil
}

// 更新题目基础信息
func (ps *ProblemService) UpdateProblem(ctx context.Context, in *pb.UpdateProblemRequest) (*emptypb.Empty, error) {
	problem := &model.Problem{
		ID:    in.Data.Id,
		Title: in.Data.Title,
		Level: int8(in.Data.Level),
	}
	var err error
	problem.Tags, err = json.Marshal(in.Data.Tags)
	if err != nil {
		logrus.Errorf("json marshal failed:%s", err.Error())
		return nil, err
	}

	if err = ps.repo.UpdateProblem(problem); err != nil {
		logrus.Errorf("UpdateProblem failed, err:%s", err.Error())
		return nil, err
	}
	return nil, nil
}

func (ps *ProblemService) DeleteProblem(ctx context.Context, in *pb.DeleteProblemRequest) (*emptypb.Empty, error) {
	return nil, ps.repo.DeleteProblem(in.Id)
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

	total, problems, err := ps.repo.QueryProblemList(int(in.Page), int(in.PageSize), in.Keyword, in.Tag, int32(role))
	if err != nil {
		return nil, err
	}
	resp.Total = total
	for _, problem := range problems {
		pbProblem := &pb.Problem{
			Id:       problem.ID,
			CreateAt: problem.CreateAt.Unix(),
			UpdateAt: problem.UpdateAt.Unix(),
			Title:    problem.Title,
			Level:    int32(problem.Level),
			Status:   int32(problem.Status),
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

	problem, err := ps.repo.QueryProblemData(in.Id)
	if err != nil {
		return nil, err
	}
	resp.Problem = &pb.Problem{
		Id:       problem.ID,
		CreateAt: problem.CreateAt.Unix(),
		UpdateAt: problem.UpdateAt.Unix(),
		Title:    problem.Title,
		Level:    int32(problem.Level),
		Status:   int32(problem.Status),
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
	problem, err := ps.repo.QueryProblemData(in.ProblemId)
	if err != nil {
		logrus.Errorf("query problem config url failed:%s", err.Error())
		return nil, err
	}

	// 获取分布式锁
	key := fmt.Sprintf("%s:%d", consts.UserLockPrefix, uid)
	_, err = ps.repo.Lock(key, consts.UserLockTTL)
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
		_ = ps.repo.UnLock(key) // 释放锁
		return nil, status.Errorf(codes.Internal, "marshal failed")
	}
	// 提交到mq
	if err = ps.judgeProducer.Publish(task_data); err != nil {
		logrus.Errorf("publish to mq failed: %v", err)
		_ = ps.repo.UnLock(key) // 释放锁
		return nil, status.Errorf(codes.Internal, "服务器错误")
	}

	return &pb.SubmitProblemResponse{
		TaskId: taskId,
	}, nil
}

func (ps *ProblemService) queryUserName(uid int64) (string, error) {
	c, cancel := context.WithTimeout(context.Background(), consts.GrpcTimeout)
	defer cancel()
	resp, err := ps.clients.User.GetUserInfoByUid(c, &pb.GetUserInfoRequest{
		Uid: uid,
	})
	if err != nil {
		return "", err
	}
	return resp.Data.Username, nil
}

func (ps *ProblemService) GetTagList(ctx context.Context, empty *emptypb.Empty) (*pb.GetTagListResponse, error) {
	resp := &pb.GetTagListResponse{}

	// 查询所有的标签
	tagList, err := ps.repo.QueryTagList()
	if err != nil {
		return nil, err
	}
	resp.Data = tagList
	return resp, nil
}

func (ps *ProblemService) ValidateProblems(ctx context.Context, in *pb.ValidateProblemsRequest) (*emptypb.Empty, error) {

}
