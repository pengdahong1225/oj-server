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
	"oj-server/module/configs"
	"oj-server/module/db"
	"oj-server/proto/pb"
	"time"
)

type JudgeRepo struct {
	db_  *gorm.DB
	rdb_ *redis.Client
}

func NewRepo() (*JudgeRepo, error) {
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

	return &JudgeRepo{
		db_:  db_,
		rdb_: rdb_,
	}, nil
}

func (r *JudgeRepo) QueryUserInfo(uid int64) (*db.UserInfo, error) {
	var user db.UserInfo
	result := r.db_.Where("id=?", uid).Find(&user)
	if result.Error != nil {
		logrus.Errorf("query user info failed, err:%v", result.Error)
		return nil, status.Errorf(codes.Internal, "query failed")
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return &user, nil
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

// 1.更新提交记录表
// 2.更新解题记录表
// 3.更新统计表
func (r *JudgeRepo) UpdateUserSubmitRecord(record *db.SubmitRecord, level int32) error {
	err := r.db_.Transaction(func(tx *gorm.DB) error {
		// 1.更新用户提交记录表
		/*
			insert into user_submit_record
			(uid, user_name, problem_id, problem_name, status, code, result, lang)
			values
			(?, ?, ?, ?, ?, ?, ?, ?);
		*/
		result := tx.Create(record)
		if result.Error != nil {
			logrus.Errorf("insert user_submit_record failed, err:%v", result.Error)
			return result.Error
		}
		// 2.更新解题记录表
		/*
			insert into user_solution(uid,problem_id)
			values(?,?)
		*/
		var repeatedAc = true
		if record.Status == "Accepted" {
			var userSolution db.UserSolution
			result = tx.Where("uid=? and problem_id=?", record.Uid, record.ProblemID).Find(&userSolution)
			if result.RowsAffected == 0 {
				repeatedAc = false
				userSolution.Uid = record.Uid
				userSolution.ProblemID = record.ProblemID
				result = tx.Create(&userSolution)
				if result.Error != nil {
					logrus.Errorf("insert user_solution failed, err:%v", result.Error)
					return result.Error
				}
			}
		}
		// 3.更新统计表
		var statistic = db.Statistics{
			Uid:    record.Uid,
			Period: time.Now().Format("2006-01"),
		}
		result = tx.FirstOrCreate(&statistic)
		if result.Error != nil {
			logrus.Errorf("query statistics failed, err:%v", result.Error)
			return result.Error
		}
		statistic.SubmitCount += 1
		if !repeatedAc && record.Status == "Accepted" {
			statistic.AccomplishCount += 1
			switch level {
			case 1:
				statistic.EasyProblemCount += 1
			case 2:
				statistic.MediumProblemCount += 1
			case 3:
				statistic.HardProblemCount += 1
			}
		}
		result = tx.Where("uid=?", statistic.Uid).Save(&statistic)
		if result.Error != nil {
			logrus.Errorf("update statistics failed, err:%v", result.Error)
			return result.Error
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (r *JudgeRepo) QueryUserAcceptCount(uid int64) (int64, int64, error) {
	var (
		dailyAcCount   int64
		monthlyAcCount int64
	)

	// 查询日 accepts
	/*
		select count(DISTINCT problem_id)
		from user_solution
		where  uid = ? and create_at DATE(create_at) = CURRENT_DATE()
		group by uid;
	*/
	result := r.db_.Table("user_solution").
		Where("uid=? and DATE(create_at) = CURRENT_DATE()", uid).
		Count(&dailyAcCount)
	if result.Error != nil {
		logrus.Errorf("query daily ac count failed, err:%v", result.Error)
		return 0, 0, result.Error
	}

	// 查询月 accepts
	/*
		select accomplish_count
		from  statistics
		where period = ? and uid = ?;
	*/
	statistic := db.Statistics{
		Uid:    uid,
		Period: time.Now().Format("2006-01"),
	}
	result = r.db_.Where("period=? and uid=?", statistic.Period, statistic.Uid).
		First(&statistic)
	if result.Error != nil {
		logrus.Errorf("query monthly ac count failed, err:%v", result.Error)
		return 0, 0, result.Error
	}
	monthlyAcCount = int64(statistic.AccomplishCount)

	return dailyAcCount, monthlyAcCount, nil
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

// lua 原子操作:
// 1. 查询当前榜尾score
// 2. 判断当前score是否比当前榜尾score大
// 3. 插入当前用户
func (r *JudgeRepo) UpdateLeaderboard(targetKey string, lb *pb.LeaderboardUserInfo) error {
	script := ``

	return nil
}
