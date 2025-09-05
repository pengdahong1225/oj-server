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
)

type JudgeRepo struct {
	db_  *gorm.DB
	rdb_ *redis.Client
}

func NewRepo() (*JudgeRepo, error) {
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

	return &JudgeRepo{
		db_:  db_,
		rdb_: rdb_,
	}, nil
}

func (r *JudgeRepo) QueryProblemData(id int64) (*db.Problem, error) {
	var problem db.Problem
	result := r.db_.Where("id=?", id).Find(&problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query failed")
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "problem not found")
	}
	return &problem, nil
}

// UpdateUserSubmitRecord
// 1.更新用户提交记录表
// 2.更新(用户,题目)AC表
// 3.更新用户解题情况统计表
func (r *JudgeRepo) UpdateUserSubmitRecord(record *db.SubmitRecord, level int32) error {
	tx := r.db_.Begin()
	if tx.Error != nil {
		logrus.Errorf("tx error: %v", tx.Error)
		return tx.Error
	}

	// todo 更新用户提交记录表
	/*
		insert into user_submit_record
		(uid, user_name, problem_id, problem_name, status, code, result, lang)
		values
		(?, ?, ?, ?, ?, ?, ?, ?);
	*/
	result := tx.Create(record)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		tx.Rollback()
		return result.Error
	}

	// todo 更新(用户,题目)AC表
	/*
		insert into user_solution(uid,problem_id)
		values(?,?)
	*/
	var repeatedAc = true
	if record.Status == "Accepted" {
		data := db.UserSolution{}
		result = tx.Where("uid=? and problem_id=?", record.Uid, record.ProblemID).Find(&data)
		if result.RowsAffected == 0 {
			repeatedAc = false
			data.Uid = record.Uid
			data.ProblemID = record.ProblemID

			result = tx.Create(&data)
			if result.Error != nil {
				logrus.Errorln(result.Error.Error())
				tx.Rollback()
				return result.Error
			}
		}
	}

	// todo 更新用户解题情况统计表
	data := db.Statistics{
		Uid: record.Uid,
	}
	result = tx.FirstOrCreate(&data)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return result.Error
	}
	data.SubmitCount += 1
	if !repeatedAc && record.Status == "Accepted" {
		data.AccomplishCount += 1
		switch level {
		case 1:
			data.EasyProblemCount += 1
		case 2:
			data.MediumProblemCount += 1
		case 3:
			data.HardProblemCount += 1
		}
	}
	result = tx.Where("uid=?", data.Uid).Save(&data)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		tx.Rollback()
		return result.Error
	}

	tx.Commit()
	return nil
}

func (r *JudgeRepo) UnLock(key string) error {
	return r.rdb_.Del(context.Background(), key).Err()
}
func (r *JudgeRepo) SetTaskState(taskId string, state int) error {
	key := fmt.Sprintf("%s:%s", global.TaskStatePrefix, taskId)
	return r.rdb_.SetEx(context.Background(), key, state, global.TaskStateExpired).Err()
}
func (r *JudgeRepo) SetTaskResult(taskId, result string) error {
	key := fmt.Sprintf("%s:%s", global.TaskResultPrefix, taskId)
	return r.rdb_.SetEx(context.Background(), key, result, global.TaskResultExpired).Err()
}
