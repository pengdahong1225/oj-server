package handler

import (
	"context"
	"db-service/global"
	"db-service/internal/models"
	pb "db-service/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mervick/aes-everywhere/go/aes256"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

var publicKey = "LeoMessi"

type DBServiceServer struct {
	pb.UnimplementedDBServiceServer
}

func (receiver *DBServiceServer) GetUserData(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var user models.UserInfo
	result := global.DBInstance.Where("phone=?", request.Phone).Find(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	if result.Error != nil {
		log.Println(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query user failed")
	}
	// 密码反解
	user.Password = aes256.Decrypt(user.Password, publicKey)
	userInfo := &pb.UserInfo{
		Phone:    user.Phone,
		Password: user.Password,
		Nickname: user.NickName,
		Email:    user.Email,
		Gender:   user.Gender,
		Role:     user.Role,
		HeadUrl:  user.HeadUrl,
	}

	return &pb.GetUserResponse{
		Data: userInfo,
	}, nil
}

func (receiver *DBServiceServer) CreateUserData(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	var user models.UserInfo
	result := global.DBInstance.Where("phone=?", request.Data.Phone)
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

	result = global.DBInstance.Create(&user)
	if result.Error != nil {
		log.Println(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "create user failed")
	}
	return &pb.CreateUserResponse{Id: user.ID}, nil
}

func (receiver *DBServiceServer) UpdateUserData(ctx context.Context, request *pb.UpdateUserRequest) (*empty.Empty, error) {
	var user models.UserInfo
	result := global.DBInstance.Where("phone=?", request.Data.Phone).Find(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	if result.Error != nil {
		log.Println(result.Error.Error())
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

	result = global.DBInstance.Save(&user) // gorm在事务执行(可重复读)，innodb自动加写锁
	if result.Error != nil {
		log.Println(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "update user failed")
	}
	return &empty.Empty{}, nil
}

func (receiver *DBServiceServer) DeleteUserData(ctx context.Context, request *pb.DeleteUserRequest) (*empty.Empty, error) {
	var user models.UserInfo
	result := global.DBInstance.Where("phone=?", request.Id).Find(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	if result.Error != nil {
		log.Println(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query user failed")
	}

	// 软删除
	result = global.DBInstance.Delete(&user)
	if result.Error != nil {
		log.Println(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "delete user failed")
	}
	// 永久删除
	// result = global.DBInstance.Unscoped().Delete(&user)

	return &empty.Empty{}, nil
}

// GetUserList 采用游标分页
func (receiver *DBServiceServer) GetUserList(ctx context.Context, request *pb.GetUserListRequest) (*pb.GetUserListResponse, error) {
	var pageSize = 10
	var userlist []models.UserInfo
	rsp := &pb.GetUserListResponse{}

	// 查询总量
	var count int64
	result := global.DBInstance.Model(&models.UserInfo{}).Count(&count)
	if result.Error != nil {
		log.Println(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query count failed")
	}
	rsp.Total = int32(count)
	if request.Cursor < 0 || int64(request.Cursor) > count {
		return nil, status.Errorf(codes.InvalidArgument, "cursor out of range")
	}
	result = global.DBInstance.Where("id >= ", request.Cursor).Order("id").Limit(pageSize).Find(&userlist)
	if result.Error != nil {
		log.Println(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query userlist failed")
	}
	for _, user := range userlist {
		rsp.Data = append(rsp.Data, &pb.UserInfo{
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
