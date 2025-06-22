package service

import (
	"context"
	"github.com/pengdahong1225/oj-server/backend/app/user/internal/respository/domain"
	"github.com/pengdahong1225/oj-server/backend/app/user/internal/respository/model"
	"github.com/pengdahong1225/oj-server/backend/module/utils"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	db *domain.MysqlDB
}

func NewUserService() *UserService {
	var err error
	s := &UserService{}
	s.db, err = domain.NewMysqlDB()
	if err != nil {
		logrus.Fatalf("NewProblemService failed, err:%s", err.Error())
	}
	return s
}

func (us *UserService) UserRegister(ctx context.Context, in *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	// 密码加密
	hash_pwd := utils.HashPassword(in.Password)

	newUser := model.UserInfo{}
}

func (us *UserService) UserLogin(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (us *UserService) ResetUserPassword(ctx context.Context, in *pb.ResetUserPasswordRequest) (*pb.ResetUserPasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (us *UserService) GetUserInfo(ctx context.Context, empty *emptypb.Empty) (*pb.GetUserInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (us *UserService) RefreshToken(ctx context.Context, empty *emptypb.Empty) (*pb.RefreshTokenResponse, error) {
	//TODO implement me
	panic("implement me")
}
