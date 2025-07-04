package domain

import (
	"github.com/sirupsen/logrus"
	"oj-server/app/common/errs"
	"oj-server/app/user/internal/respository/model"
)

func (md *MysqlDB) CreateNewUser(user *model.UserInfo) (int64, error) {
	var u model.UserInfo
	result := md.db_.Where("mobile=?", user.Mobile).Find(&u)
	if result.Error != nil {
		logrus.Errorf("查询错误, err: %s", result.Error.Error())
		return -1, errs.QueryFailed
	}
	if result.RowsAffected > 0 {
		return -1, errs.AlreadyExists
	}

	result = md.db_.Create(&user)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return -1, errs.CreateFailed
	}
	return user.ID, nil
}

func (md *MysqlDB) GetUserInfoByMobile(mobile int64) (*model.UserInfo, error) {
	var user model.UserInfo
	result := md.db_.Where("mobile=?", mobile).Find(&user)
	if result.Error != nil {
		logrus.Errorf("查询错误, err: %s", result.Error.Error())
		return nil, errs.QueryFailed
	}
	if result.RowsAffected == 0 {
		return nil, errs.NotFound
	}
	return &user, nil
}
