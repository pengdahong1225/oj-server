package domain

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"oj-server/module/model"
)

func (md *MysqlDB) CreateNewUser(user *model.UserInfo) (int64, error) {
	var u model.UserInfo
	result := md.db_.Where("mobile=?", user.Mobile).Find(&u)
	if result.Error != nil {
		logrus.Errorf("查询错误, err: %s", result.Error.Error())
		return -1, status.Errorf(codes.Internal, "query user faild")
	}
	if result.RowsAffected > 0 {
		return -1, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	result = md.db_.Create(&user)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return -1, status.Errorf(codes.Internal, "create user faild")
	}
	return user.ID, nil
}
func (md *MysqlDB) GetUserInfoByMobile(mobile int64) (*model.UserInfo, error) {
	var user model.UserInfo
	result := md.db_.Where("mobile=?", mobile).Find(&user)
	if result.Error != nil {
		logrus.Errorf("查询错误, err: %s", result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query user faild")
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return &user, nil
}

func (md *MysqlDB) ResetUserPassword(mobile int64, password string) error {
	/*
		update user_info set password = '123456'
		where mobile = ?;
	*/
	result := md.db_.Model(&model.UserInfo{}).Where("mobile=?", mobile).Update("password", password)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return status.Errorf(codes.Internal, "update user faild")
	}
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "user not found")
	}
	return nil
}
