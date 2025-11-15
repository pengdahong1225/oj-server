package service

import (
	"context"
	"strconv"

	"oj-server/svr/user/internal/biz"
	"oj-server/svr/user/internal/data"
	"oj-server/utils"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"oj-server/pkg/proto/pb"
	"oj-server/svr/user/internal/model"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	uc *biz.UserUseCase
}

func NewUserService() *UserService {
	var err error
	s := &UserService{}
	up, err := data.NewUserRepo()
	if err != nil {
		logrus.Fatalf("NewUserService failed, err:%s", err.Error())
	}
	s.uc = biz.NewUserUseCase(up) // 注入实现
	return s
}

// 注册
func (us *UserService) UserRegister(ctx context.Context, in *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	resp := &pb.UserRegisterResponse{}

	// 不可逆加密
	hash_pwd := utils.HashPassword(in.Password)
	mobile, err := strconv.ParseInt(in.Mobile, 10, 64)
	if err != nil {
		logrus.Errorf("parse mobile failed: %s", err.Error())
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
	id, err = us.uc.CreateNewUser(&newUser)
	if err != nil {
		logrus.Infof("创建新用户失败, err: %s", err.Error())
		return resp, err
	}
	resp.Uid = id
	return resp, nil
}

// 账号密码登录
func (us *UserService) UserLogin(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	// 拉取用户信息
	mobile, _ := strconv.ParseInt(in.Mobile, 10, 64)
	userInfo, err := us.uc.GetUserInfoByMobile(mobile)
	if err != nil {
		return nil, err
	}
	// 校验密码
	if utils.HashPassword(in.Password) != userInfo.PassWord {
		return nil, status.Errorf(codes.Unauthenticated, "密码错误")
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

// 验证码登录
func (us *UserService) UserLoginBySmsCode(ctx context.Context, in *pb.UserLoginBySmsCodeRequest) (*pb.UserLoginResponse, error) {
	// 拉取用户信息
	mobile, _ := strconv.ParseInt(in.Mobile, 10, 64)
	userInfo, err := us.uc.GetUserInfoByMobile(mobile)
	if err != nil {
		return nil, err
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

// 重置密码
func (us *UserService) ResetUserPassword(ctx context.Context, in *pb.ResetUserPasswordRequest) (*emptypb.Empty, error) {
	hash := utils.HashPassword(in.Password)
	mobile, _ := strconv.ParseInt(in.Mobile, 10, 64)
	err := us.uc.ResetUserPassword(mobile, hash)
	if err != nil {
		logrus.Errorf("重置密码失败, err: %s", err.Error())
		return nil, err
	}
	return nil, nil
}

// 获取用户信息
func (us *UserService) GetUserInfoByUid(ctx context.Context, in *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	// 拉取用户信息
	userInfo, err := us.uc.GetUserInfoByUid(in.Uid)
	if err != nil {
		logrus.Errorf("获取用户信息失败, err: %s", err.Error())
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

// 获取用户信息
func (us *UserService) GetUserInfoByMobile(ctx context.Context, in *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	// 拉取用户信息
	mobile, _ := strconv.ParseInt(in.Mobile, 10, 64)
	userInfo, err := us.uc.GetUserInfoByMobile(mobile)
	if err != nil {
		logrus.Errorf("获取用户信息失败, err: %s", err.Error())
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
