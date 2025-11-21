package data

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"oj-server/global"
	"oj-server/svr/user/internal/configs"
	"oj-server/svr/user/internal/model"
	"os"
	"time"
)

type UserRepo struct {
	db_  *gorm.DB
	rdb_ *redis.Client
}

func NewUserRepo() (*UserRepo, error) {
	// mysql
	mysql_cfg := configs.AppConf.MysqlCfg
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysql_cfg.User,
		mysql_cfg.Pwd, mysql_cfg.Host, mysql_cfg.Port, mysql_cfg.Db)
	timer := time.Now().Format("2006_01_02")
	filePath := fmt.Sprintf("%s/orm.%s.log", global.LogPath, timer)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	writer := io.MultiWriter(os.Stdout, file)
	newLogger := logger.New(
		log.New(writer, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)
	db_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		// SkipDefaultTransaction: true, //全局禁用默认事务
	})
	if err != nil {
		return nil, err
	}

	// redis
	redis_cfg := configs.AppConf.RedisCfg
	dsn = fmt.Sprintf("%s:%d", redis_cfg.Host, redis_cfg.Port)
	rdb_ := redis.NewClient(&redis.Options{
		Addr:    dsn,
		Network: "tcp",
	})
	st := rdb_.Ping(context.Background())
	if st.Err() != nil {
		return nil, st.Err()
	}

	return &UserRepo{
		db_:  db_,
		rdb_: rdb_,
	}, nil
}

func (up *UserRepo) CreateNewUser(user *model.UserInfo) (int64, error) {
	var u model.UserInfo
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
func (up *UserRepo) GetUserInfoByUid(mobile int64) (*model.UserInfo, error) {
	var user model.UserInfo
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
func (up *UserRepo) GetUserInfoByMobile(mobile int64) (*model.UserInfo, error) {
	var user model.UserInfo
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
	result := up.db_.Model(&model.UserInfo{}).Where("mobile=?", mobile).Update("password", password)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return status.Errorf(codes.Internal, "update user faild")
	}
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "user not found")
	}
	return nil
}

// 分页查询用户列表
// 查询{id，mobile，nickname，create_at, role}
// @page 页码
// @page_size 单页数量
// @keyword 关键字(昵称)
func (up *UserRepo) QueryUserList(page, pageSize int, keyword string) (int64, []model.UserInfo, error) {
	var count int64
	query := up.db_.Model(&model.UserInfo{})
	if keyword != "" {
		query = query.Where("nickname LIKE ?", "%"+keyword+"%")
	}
	result := query.Count(&count)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return -1, nil, status.Errorf(codes.Internal, "query failed")
	}
	offSet := (page - 1) * pageSize
	var users []model.UserInfo
	result = query.Select("id, mobile, nickname, email, create_at, role").Offset(offSet).Limit(pageSize).Find(&users)
	if result.Error != nil {
		logrus.Errorf("query user list failed: %s", result.Error.Error())
		return -1, nil, status.Errorf(codes.Internal, "query failed")
	}
	return count, users, nil
}
