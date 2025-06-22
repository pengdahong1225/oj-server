package service

import (
	"context"
	"fmt"
	"github.com/pengdahong1225/oj-server/backend/app/user/internal/respository/domain"
	"github.com/pengdahong1225/oj-server/backend/app/user/internal/respository/model"
	"github.com/pengdahong1225/oj-server/backend/module/utils"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
)

const (
	compareErr      = "密码校验失败"
	loginSucc       = "登录成功"
	getUserInfoFail = "获取用户信息失败: "
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
	resp := &pb.UserRegisterResponse{}

	// 不可逆加密
	hash_pwd := utils.HashPassword(in.Password)
	mobile, err := strconv.ParseInt(in.Mobile, 10, 64)
	if err != nil {
		logrus.Errorf("mobile转换失败, err: %s", err.Error())
		resp.Message = "mobile转换失败: " + err.Error()
		return resp, err
	}
	newUser := model.UserInfo{
		Mobile:    mobile,
		PassWord:  hash_pwd,
		NickName:  in.Nickname,
		Email:     in.Email,
		Gender:    in.Gender,
		Role:      in.Role,
		AvatarUrl: "",
	}

	var id int64
	id, err = us.db.CreateNewUser(&newUser)
	if err != nil {
		logrus.Infof("创建新用户失败, err: %s", err.Error())
		resp.Message = "创建新用户失败: " + err.Error()
		return resp, err
	}
	resp.Uid = id
	resp.Message = "创建新用户成功"
	return resp, nil
}

func (us *UserService) UserLogin(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	// 拉取用户信息
	mobile, _ := strconv.ParseInt(in.Mobile, 10, 64)
	userInfo, err := us.db.GetUserInfoByMobile(mobile)
	if err != nil {
		return nil, err
	}
	// 校验密码
	if utils.HashPassword(in.Password) != userInfo.PassWord {
		return nil, fmt.Errorf("密码错误")
	}

	resp := &pb.UserLoginResponse{
		Uid:       userInfo.ID,
		Mobile:    strconv.FormatInt(userInfo.Mobile, 10),
		NickName:  userInfo.NickName,
		Email:     userInfo.Email,
		Gender:    userInfo.Gender,
		Role:      userInfo.Role,
		AvatarUrl: userInfo.AvatarUrl,
	}
	return resp, nil
}

func (us *UserService) ResetUserPassword(ctx context.Context, in *pb.ResetUserPasswordRequest) (*pb.ResetUserPasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (us *UserService) GetUserInfo(ctx context.Context, in *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	// 拉取用户信息
	mobile, _ := strconv.ParseInt(in.Mobile, 10, 64)
	userInfo, err := us.db.GetUserInfoByMobile(mobile)
	if err != nil {
		return nil, err
	}
	u := pb.UserInfo{
		Uid:       userInfo.ID,
		CreateAt:  userInfo.CreateAt.Unix(),
		Mobile:    userInfo.Mobile,
		Nickname:  userInfo.NickName,
		Email:     userInfo.Email,
		Gender:    userInfo.Gender,
		Role:      userInfo.Role,
		AvatarUrl: userInfo.AvatarUrl,
	}
	resp := &pb.GetUserInfoResponse{
		Data: &u,
	}

	return resp, nil
}

func (us *UserService) RefreshToken(ctx context.Context, empty *emptypb.Empty) (*pb.RefreshTokenResponse, error) {
	//TODO implement me
	panic("implement me")
}
