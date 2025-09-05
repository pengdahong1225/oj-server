package service

import (
	"context"
	"fmt"
	"strconv"

	"oj-server/svr/user/internal/biz"
	"oj-server/svr/user/internal/data"
	"oj-server/utils"

	"github.com/sirupsen/logrus"
	"oj-server/module/db"
	"oj-server/module/proto/pb"
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
	if err != nil {
		logrus.Fatalf("NewUserService failed, err:%s", err.Error())
	}
	return s
}

// 注册
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
	newUser := db.UserInfo{
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
		resp.Message = "创建新用户失败: " + err.Error()
		return resp, err
	}
	resp.Uid = id
	resp.Message = "创建新用户成功"
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
func (us *UserService) ResetUserPassword(ctx context.Context, in *pb.ResetUserPasswordRequest) (*pb.ResetUserPasswordResponse, error) {
	resp := &pb.ResetUserPasswordResponse{}
	hash := utils.HashPassword(in.Password)
	mobile, _ := strconv.ParseInt(in.Mobile, 10, 64)
	err := us.uc.ResetUserPassword(mobile, hash)
	if err != nil {
		logrus.Errorf("重置密码失败, err: %s", err.Error())
		resp.Message = "重置密码失败: " + err.Error()
		return resp, err
	}
	return resp, nil
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
