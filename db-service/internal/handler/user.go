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
		return nil, result.Error
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
		HeadPic:  user.HeadUrl,
	}

	return &pb.GetUserResponse{
		Data: userInfo,
	}, nil
}

func (receiver *DBServiceServer) CreateUserData(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	return nil, status.Errorf(codes.Unimplemented, "method CreateUserData not implemented")
}

func (receiver *DBServiceServer) UpdateUserData(ctx context.Context, request *pb.UpdateUserRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserData not implemented")
}

func (receiver *DBServiceServer) DeleteUserData(ctx context.Context, request *pb.DeleteUserRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserData not implemented")
}

func (receiver *DBServiceServer) GetUserList(ctx context.Context, request *pb.GetUserListRequest) (*pb.GetUserListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserList not implemented")
}
