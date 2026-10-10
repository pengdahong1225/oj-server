package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"oj-server/common/proto/pb"
	"oj-server/svr/user/internal/client"
	"oj-server/svr/user/internal/data"
	"strconv"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	repo    *data.Repo
	clients *client.Clients
}

func NewUserService(clients *client.Clients) (*UserService, error) {
	repo, err := data.NewRepo()
	if err != nil {
		return nil, err
	}

	return &UserService{
		repo:    repo,
		clients: clients,
	}, nil
}

// 注册
func (us *UserService) UserRegister(ctx context.Context, in *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	resp := &pb.UserRegisterResponse{}

	// 不可逆加密
	hash_pwd := hashPassword(in.Password)
	mobile, err := strconv.ParseInt(in.Mobile, 10, 64)
	if err != nil {
		logrus.Errorf("parse mobile failed: %s", err.Error())
		return resp, err
	}
	newUser := data.UserInfo{
		Mobile:    mobile,
		PassWord:  hash_pwd,
		NickName:  in.Nickname,
		Email:     in.Email,
		Gender:    in.Gender,
		Role:      in.Role,
		AvatarUrl: "",
	}

	var id int64
	id, err = us.repo.CreateNewUser(&newUser)
	if err != nil {
		logrus.Infof("创建新用户失败, err: %s", err.Error())
		return resp, err
	}
	resp.Uid = id
	return resp, nil
}

func hashPassword(from string) string {
	hash := sha256.Sum256([]byte(from))
	return hex.EncodeToString(hash[:])
}

// 账号密码登录
func (us *UserService) UserLogin(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	mobile, err := strconv.ParseInt(in.Mobile, 10, 64)
	if err != nil {
		logrus.Errorf("parse mobile failed: %s", err.Error())
		return nil, err
	}

	userInfo, err := us.repo.GetUserInfoByMobile(mobile)
	if err != nil {
		logrus.Errorf("get user by mobile err: %s", err.Error())
		return nil, err
	}
	if hashPassword(in.Password) != userInfo.PassWord {
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
	userInfo, err := us.repo.GetUserInfoByMobile(mobile)
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
	hash := hashPassword(in.Password)
	mobile, _ := strconv.ParseInt(in.Mobile, 10, 64)
	err := us.repo.ResetUserPassword(mobile, hash)
	if err != nil {
		logrus.Errorf("重置密码失败, err: %s", err.Error())
		return nil, err
	}
	return nil, nil
}

// 获取用户信息
func (us *UserService) GetUserInfoByUid(ctx context.Context, in *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	// 拉取用户信息
	userInfo, err := us.repo.GetUserInfoByUid(in.Uid)
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

func (us *UserService) GetUserList(ctx context.Context, in *pb.GetUserListRequest) (*pb.GetUserListResponse, error) {
	resp := &pb.GetUserListResponse{}
	total, users, err := us.repo.QueryUserList(int(in.Page), int(in.PageSize), in.Keyword)
	if err != nil {
		return nil, err
	}
	resp.Total = total

	for _, user := range users {
		u := pb.UserInfo{
			Uid:      user.ID,
			CreateAt: user.CreateAt.Unix(),
			Mobile:   user.Mobile,
			Nickname: user.NickName,
			Email:    user.Email,
			Role:     user.Role,
		}
		resp.Data = append(resp.Data, &u)
	}

	return resp, nil
}
