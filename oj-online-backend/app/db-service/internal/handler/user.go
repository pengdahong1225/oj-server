package handler

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/services/mysql"
	"github.com/pengdahong1225/Oj-Online-Server/proto/pb"
	"github.com/sirupsen/logrus"
)

func (receiver *DBServiceServer) GetUserDataByMobile(ctx context.Context, request *pb.GetUserDataByMobileRequest) (*pb.GetUserResponse, error) {
	db := mysql.Instance()
	var user mysql.UserInfo
	result := db.Where("mobile=?", request.Mobile).Find(&user)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, QueryField
	}
	if result.RowsAffected == 0 {
		return nil, NotFound
	}

	userInfo := &pb.UserInfo{
		Uid:       user.ID,
		Mobile:    user.Mobile,
		Nickname:  user.NickName,
		Email:     user.Email,
		Gender:    user.Gender,
		Role:      user.Role,
		AvatarUrl: user.AvatarUrl,
	}

	return &pb.GetUserResponse{
		Data: userInfo,
	}, nil
}

func (receiver *DBServiceServer) GetUserDataByUid(ctx context.Context, request *pb.GetUserDataByUidRequest) (*pb.GetUserResponse, error) {
	db := mysql.Instance()
	var user mysql.UserInfo
	result := db.Where("id=?", request.Id).Find(&user)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, QueryField
	}
	if result.RowsAffected == 0 {
		return nil, NotFound
	}

	userInfo := &pb.UserInfo{
		Uid:       user.ID,
		Mobile:    user.Mobile,
		Nickname:  user.NickName,
		Email:     user.Email,
		Gender:    user.Gender,
		Role:      user.Role,
		AvatarUrl: user.AvatarUrl,
	}

	return &pb.GetUserResponse{
		Data: userInfo,
	}, nil
}

func (receiver *DBServiceServer) CreateUserData(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	db := mysql.Instance()
	var user mysql.UserInfo
	result := db.Where("mobile=?", request.Data.Mobile)
	if result.RowsAffected > 0 {
		return nil, AlreadyExists
	}

	user.Mobile = request.Data.Mobile
	user.NickName = request.Data.Nickname
	user.Email = request.Data.Email
	user.Gender = request.Data.Gender
	user.Role = request.Data.Role
	user.AvatarUrl = request.Data.AvatarUrl

	result = db.Create(&user)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, InsertFailed
	}

	l, _ := json.Marshal(user)
	logrus.Debugln("create user:%s", l)

	return &pb.CreateUserResponse{Id: user.ID}, nil
}

func (receiver *DBServiceServer) UpdateUserData(ctx context.Context, request *pb.UpdateUserRequest) (*empty.Empty, error) {
	db := mysql.Instance()
	var user mysql.UserInfo
	result := db.Where("mobile=?", request.Data.Mobile).Find(&user)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, QueryField
	}
	if result.RowsAffected == 0 {
		return nil, NotFound
	}

	// 更新
	user.NickName = request.Data.Nickname
	user.Email = request.Data.Email
	user.Gender = request.Data.Gender
	user.Role = request.Data.Role
	user.AvatarUrl = request.Data.AvatarUrl

	result = db.Save(&user) // gorm在事务执行(可重复读)，innodb自动加写锁
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, InsertFailed
	}
	return &empty.Empty{}, nil
}

func (receiver *DBServiceServer) DeleteUserData(ctx context.Context, request *pb.DeleteUserRequest) (*empty.Empty, error) {
	db := mysql.Instance()
	var user mysql.UserInfo
	result := db.Where("id=?", request.Id).Find(&user)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, QueryField
	}
	if result.RowsAffected == 0 {
		return nil, NotFound
	}

	// 软删除
	result = db.Delete(&user)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, DeleteFailed
	}
	// 永久删除
	// result = global.DBInstance.Unscoped().Delete(&user)

	return &empty.Empty{}, nil
}

// GetUserList 采用游标分页
func (receiver *DBServiceServer) GetUserList(ctx context.Context, request *pb.GetUserListRequest) (*pb.GetUserListResponse, error) {
	// db := mysql.Instance()
	// var pageSize = 10
	// var userlist []models.UserInfo
	// rsp := &pb.GetUserListResponse{}
	//
	// // 查询总量
	// var count int64
	// result := db.Model(&models.UserInfo{}).Count(&count)
	// if result.Error != nil {
	//	logrus.Debugln(result.Error.Error())
	//	return nil, status.Errorf(codes.Internal, "query count failed")
	// }
	// rsp.Total = int32(count)
	// if request.Cursor < 0 || int64(request.Cursor) > count {
	//	return nil, status.Errorf(codes.InvalidArgument, "cursor out of range")
	// }
	// result = db.Where("id >= ", request.Cursor).Order("id").Limit(pageSize).Find(&userlist)
	// if result.Error != nil {
	//	logrus.Debugln(result.Error.Error())
	//	return nil, status.Errorf(codes.Internal, "query userlist failed")
	// }
	// for _, user := range userlist {
	//	rsp.Data = append(rsp.Data, &pb.UserInfo{
	//		Phone:    user.Phone,
	//		Nickname: user.NickName,
	//		Email:    user.Email,
	//		Gender:   user.Gender,
	//		Role:     user.Role,
	//		HeadUrl:  user.HeadUrl,
	//	})
	// }
	// rsp.Cursor = request.Cursor + int32(result.RowsAffected) + 1

	return nil, nil
}

// GetUserSolvedList 查询用户哪些题目
func (receiver *DBServiceServer) GetUserSolvedList(ctx context.Context, request *pb.GetUserSolvedListRequest) (*pb.GetUserSolvedListResponse, error) {
	db := mysql.Instance()
	var userSolutionList []mysql.UserSolution
	result := db.Where("uid=?", request.Uid).Find(&userSolutionList)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, QueryField
	}
	if result.RowsAffected == 0 {
		return nil, NotFound
	}

	var rsp = new(pb.GetUserSolvedListResponse)
	for _, userSolution := range userSolutionList {
		rsp.ProblemSolvedList = append(rsp.ProblemSolvedList, userSolution.ProblemID)
	}

	return rsp, nil
}
