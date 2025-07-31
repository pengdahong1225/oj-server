package data

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"oj-server/global"
	"oj-server/module/configManager"
	"oj-server/module/db"
	"oj-server/module/db/model"
	"time"
)

type RecordRepo struct {
	db_  *gorm.DB
	rdb_ *redis.Client
}

func NewRecordRepo() (*RecordRepo, error) {
	mysql_cfg := configManager.AppConf.MysqlCfg
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysql_cfg.User,
		mysql_cfg.Pwd, mysql_cfg.Host, mysql_cfg.Port, mysql_cfg.Db)
	db_, err := db.NewMysqlCli(dsn, global.LogPath)
	if err != nil {
		return nil, err
	}

	redis_cfg := configManager.AppConf.RedisCfg
	dsn = fmt.Sprintf("%s:%d", redis_cfg.Host, redis_cfg.Port)
	rdb_, err := db.NewRedisCli(dsn)
	if err != nil {
		return nil, err
	}

	return &RecordRepo{
		db_:  db_,
		rdb_: rdb_,
	}, nil
}

// 查询用户的历史提交记录
func (rr *RecordRepo) QuerySubmitRecordList(uid int64, pageSize, offset int) (int64, []model.SubmitRecord, error) {
	/*
		select id, created_at, problem_name, status, lang from user_submit_record
		where uid = ?
		order by created_at desc
		offset off_set
		limit page_size;
	*/
	var count int64 = 0
	result := rr.db_.Model(&model.SubmitRecord{}).Where("uid = ?", uid).Count(&count)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return 0, nil, status.Errorf(codes.Internal, "查询提交记录失败")
	}
	var records []model.SubmitRecord
	result = rr.db_.Where("uid = ?", uid).Order("created_at desc").Offset(offset).Limit(pageSize).Find(&records)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return 0, nil, status.Errorf(codes.Internal, "查询提交记录失败")
	}
	return count, records, nil
}

// 查询某项提交记录的详细信息
func (rr *RecordRepo) QuerySubmitRecord(id int64) (*model.SubmitRecord, error) {
	var record model.SubmitRecord
	result := rr.db_.Where("id = ?", id).First(&record)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query failed")
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "record not found")
	}
	return &record, nil
}

func (rr *RecordRepo) Lock(key string, ttl time.Duration) (bool, error) {
	return rr.rdb_.SetNX(context.Background(), key, "locked", ttl).Result()
}
func (rr *RecordRepo) UnLock(key string) error {
	return rr.rdb_.Del(context.Background(), key).Err()
}
