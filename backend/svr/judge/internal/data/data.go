package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"oj-server/global"

	"oj-server/pkg/proto/pb"
	"oj-server/svr/judge/internal/configs"
	"oj-server/svr/problem/internal/model"
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

func (r *JudgeRepo) QueryUserInfo(uid int64) (*model.UserInfo, error) {
	var user model.UserInfo
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

func (r *JudgeRepo) QueryProblemData(id int64) (*model2.Problem, error) {
	var problem model2.Problem
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
func (r *JudgeRepo) UpdateUserSubmitRecord(record *model2.SubmitRecord, level int32) error {
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
			var userSolution model2.UserSolution
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
		var statistic = model2.Statistics{
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
	statistic := model2.Statistics{
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
	script := `
-- KEYS[1]: 排行榜的key
-- ARGV[1]: 用户JSON信息
-- ARGV[2]: 最大显示数量(200)
-- ARGV[3]: 最大容错数量(100)

-- 解析参数
local maxDisplay = tonumber(ARGV[2])
local maxCapacity = maxDisplay + tonumber(ARGV[3])
local userData = cjson.decode(ARGV[1])
local userId = userData.uid
local newScore = tonumber(userData.score)

-- 获取当前排行榜信息
local currentCount = redis.call('ZCARD', KEYS[1])
local currentScore = redis.call('ZSCORE', KEYS[1], userId)

-- 判断用户是否能进入排行榜
local canEnter = false

-- 情况1: 排行榜未满(小于最大容量)
if currentCount < maxCapacity then
    canEnter = true
-- 情况2: 用户已在排行榜中
elseif currentScore ~= false then
    canEnter = true
-- 情况3: 新分数高于排行榜最低分
else
    local lowestScore = redis.call('ZRANGE', KEYS[1], -1, -1, 'WITHSCORES')[2]
    if newScore > tonumber(lowestScore) then
        canEnter = true
    end
end

-- 只有符合条件的用户才能更新
if canEnter then
    -- 更新用户分数和信息
    redis.call('ZADD', KEYS[1], newScore, userId)
    redis.call('HSET', KEYS[1] .. ':user_data', userId, ARGV[1])
    
    -- 维护排行榜大小(不超过最大容量)
    local currentCount = redis.call('ZCARD', KEYS[1])
    if currentCount > maxCapacity then
        -- 移除超出容量的最低分用户
        redis.call('ZREMRANGEBYRANK', KEYS[1], 0, currentCount - maxCapacity - 1)
    end
end

-- 强制清理超出显示数量的数据(保持严格限制)
local displayCount = redis.call('ZCARD', KEYS[1])
if displayCount > maxDisplay then
    redis.call('ZREMRANGEBYRANK', KEYS[1], 0, displayCount - maxDisplay - 1)
end
`
	// 序列化用户信息
	userJson, err := json.Marshal(lb)
	if err != nil {
		return fmt.Errorf("failed to marshal user info: %v", err)
	}

	// 执行脚本
	_, err = r.rdb_.Eval(
		context.Background(),
		script,
		[]string{targetKey},
		string(userJson),
		200, // 最大显示数量
		100, // 容错数量
	).Result()

	return nil
}
