package handler

import (
	"context"
	"db-service/internal/models"
	pb2 "db-service/internal/proto"
	"db-service/services/dao/mysql"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mervick/aes-everywhere/go/aes256"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var publicKey = "LeoMessi"

type DBServiceServer struct {
	pb2.UnimplementedDBServiceServer
}

func (receiver *DBServiceServer) GetUserData(ctx context.Context, request *pb2.GetUserRequest) (*pb2.GetUserResponse, error) {
	db := mysql.DB
	var user models.UserInfo
	result := db.Where("phone=?", request.Phone).Find(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query user failed")
	}
	// 密码反解
	user.Password = aes256.Decrypt(user.Password, publicKey)
	userInfo := &pb2.UserInfo{
		Phone:    user.Phone,
		Password: user.Password,
		Nickname: user.NickName,
		Email:    user.Email,
		Gender:   user.Gender,
		Role:     user.Role,
		HeadUrl:  user.HeadUrl,
	}

	return &pb2.GetUserResponse{
		Data: userInfo,
	}, nil
}

func (receiver *DBServiceServer) CreateUserData(ctx context.Context, request *pb2.CreateUserRequest) (*pb2.CreateUserResponse, error) {
	db := mysql.DB
	var user models.UserInfo
	result := db.Where("phone=?", request.Data.Phone)
	if result.RowsAffected > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}
	user.Phone = request.Data.Phone
	user.Password = aes256.Encrypt(request.Data.Password, publicKey)
	user.NickName = request.Data.Nickname
	user.Email = request.Data.Email
	user.Gender = request.Data.Gender
	user.Role = request.Data.Role
	user.HeadUrl = request.Data.HeadUrl

	l, _ := json.Marshal(user)
	logrus.Infof("create user:%s", l)

	result = db.Create(&user)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "create user failed")
	}
	return &pb2.CreateUserResponse{Id: user.ID}, nil
}

func (receiver *DBServiceServer) UpdateUserData(ctx context.Context, request *pb2.UpdateUserRequest) (*empty.Empty, error) {
	db := mysql.DB
	var user models.UserInfo
	result := db.Where("phone=?", request.Data.Phone).Find(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query user failed")
	}
	// 密码反解
	user.Password = aes256.Decrypt(user.Password, publicKey)

	// 更新
	user.NickName = request.Data.Nickname
	user.Password = request.Data.Password
	user.Role = request.Data.Role
	user.Gender = request.Data.Gender
	user.Email = request.Data.Email
	user.HeadUrl = request.Data.HeadUrl

	result = db.Save(&user) // gorm在事务执行(可重复读)，innodb自动加写锁
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "update user failed")
	}
	return &empty.Empty{}, nil
}

func (receiver *DBServiceServer) DeleteUserData(ctx context.Context, request *pb2.DeleteUserRequest) (*empty.Empty, error) {
	db := mysql.DB
	var user models.UserInfo
	result := db.Where("phone=?", request.Id).Find(&user)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query user failed")
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	// 软删除
	result = db.Delete(&user)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "delete user failed")
	}
	// 永久删除
	// result = global.DBInstance.Unscoped().Delete(&user)

	return &empty.Empty{}, nil
}

// GetUserList 采用游标分页
func (receiver *DBServiceServer) GetUserList(ctx context.Context, request *pb2.GetUserListRequest) (*pb2.GetUserListResponse, error) {
	db := mysql.DB
	var pageSize = 10
	var userlist []models.UserInfo
	rsp := &pb2.GetUserListResponse{}

	// 查询总量
	var count int64
	result := db.Model(&models.UserInfo{}).Count(&count)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query count failed")
	}
	rsp.Total = int32(count)
	if request.Cursor < 0 || int64(request.Cursor) > count {
		return nil, status.Errorf(codes.InvalidArgument, "cursor out of range")
	}
	result = db.Where("id >= ", request.Cursor).Order("id").Limit(pageSize).Find(&userlist)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query userlist failed")
	}
	for _, user := range userlist {
		rsp.Data = append(rsp.Data, &pb2.UserInfo{
			Phone:    user.Phone,
			Nickname: user.NickName,
			Email:    user.Email,
			Gender:   user.Gender,
			Role:     user.Role,
			HeadUrl:  user.HeadUrl,
		})
	}
	rsp.Cursor = request.Cursor + int32(result.RowsAffected) + 1

	return rsp, nil
}
