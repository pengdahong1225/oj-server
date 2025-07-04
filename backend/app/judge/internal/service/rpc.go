package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"oj-server/app/judge/internal/respository/domain"
	"oj-server/proto/pb"
)

type JudgeService struct {
	pb.UnimplementedJudgeServiceServer
	db *domain.MysqlDB
}

func NewJudgeService() *JudgeService {
	var err error
	s := &JudgeService{}
	s.db, err = domain.NewMysqlDB()
	if err != nil {
		logrus.Fatalf("NewJudgeService failed, err:%s", err.Error())
	}
	return s
}

func (js *JudgeService) QueryJudgeResult(ctx context.Context, in *pb.QueryJudgeResultRequest) (*pb.QueryJudgeResultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryJudgeResult not implemented")
}
