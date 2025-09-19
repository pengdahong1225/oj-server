package data

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"oj-server/global"
	"oj-server/module/configs"
	"oj-server/module/db"
)

type UserRepo struct {
	db_  *gorm.DB
	rdb_ *redis.Client
}

func NewUserRepo() (*UserRepo, error) {
	mysql_cfg := configs.AppConf.MysqlCfg
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysql_cfg.User,
		mysql_cfg.Pwd, mysql_cfg.Host, mysql_cfg.Port, mysql_cfg.Db)
	db_, err := db.NewMysqlCli(dsn, global.LogPath)
	if err != nil {
		return nil, err
	}

	redis_cfg := configs.AppConf.RedisCfg
	dsn = fmt.Sprintf("%s:%d", redis_cfg.Host, redis_cfg.Port)
	rdb_, err := db.NewRedisCli(dsn)
	if err != nil {
		return nil, err
	}

	return &UserRepo{
		db_:  db_,
		rdb_: rdb_,
	}, nil
}

func (up *UserRepo) CreateNewUser(user *db.UserInfo) (int64, error) {
	var u db.UserInfo
	result := up.db_.Where("mobile=?", user.Mobile).Find(&u)
	if result.Error != nil {
		logrus.Errorf("查询错误, err: %s", result.Error.Error())
		return -1, status.Errorf(codes.Internal, "query user faild")
	}
	if result.RowsAffected > 0 {
		return -1, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	result = up.db_.Create(&user)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return -1, status.Errorf(codes.Internal, "create user faild")
	}
	return user.ID, nil
}
func (up *UserRepo) GetUserInfoByUid(mobile int64) (*db.UserInfo, error) {
	var user db.UserInfo
	result := up.db_.Where("id=?", mobile).Find(&user)
	if result.Error != nil {
		logrus.Errorf("查询错误, err: %s", result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query user faild")
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return &user, nil
}
func (up *UserRepo) GetUserInfoByMobile(mobile int64) (*db.UserInfo, error) {
	var user db.UserInfo
	result := up.db_.Where("mobile=?", mobile).Find(&user)
	if result.Error != nil {
		logrus.Errorf("查询错误, err: %s", result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query user faild")
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return &user, nil
}

func (up *UserRepo) ResetUserPassword(mobile int64, password string) error {
	/*
		update user_info set password = '123456'
		where mobile = ?;
	*/
	result := up.db_.Model(&db.UserInfo{}).Where("mobile=?", mobile).Update("password", password)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return status.Errorf(codes.Internal, "update user faild")
	}
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "user not found")
	}
	return nil
}
