package data

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"oj-server/global"
	"oj-server/pkg/proto/pb"
	"oj-server/svr/problem/internal/configs"
	"oj-server/svr/problem/internal/model"
	"os"
	"strconv"
	"time"
)

type RecordRepo struct {
	db_  *gorm.DB
	rdb_ *redis.Client
}

func NewRecordRepo() (*RecordRepo, error) {
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

func (rr *RecordRepo) QueryLeaderboardLastUpdate() (int64, error) {
	result := rr.rdb_.Get(context.Background(), global.LeaderboardLastUpdateKey)
	if result.Err() != nil {
		return 0, result.Err()
	}
	return result.Int64()
}
func (rr *RecordRepo) UpdateLeaderboardLastUpdate(time int64) error {
	result := rr.rdb_.Set(context.Background(), global.LeaderboardLastUpdateKey, time, redis.KeepTTL)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}
func (rr *RecordRepo) QueryStatistics(uid int64) (*model.Statistics, error) {
	var statistics model.Statistics
	result := rr.db_.Where("uid = ?", uid).First(&statistics)
	if result.Error != nil {
		return nil, result.Error
	}
	return &statistics, nil
}

// 排行榜临时数据
type leaderboardData struct {
	Uid             int64 `gorm:"column:uid"`
	AccomplishCount int32 `gorm:"column:accomplish_count"`

	Mobile   int64  `gorm:"column:mobile"`
	Username string `gorm:"column:nickname"`
	Avatar   string `gorm:"column:avatar_url"`
}

func (rr *RecordRepo) QueryMonthAccomplishLeaderboard(limit int, period string) ([]*pb.LeaderboardUserInfo, error) {
	var lb_datas []leaderboardData
	/*
		SELECT
		    s.uid, s.accomplish_count,
		    u.nickname AS username, u.avatar_url AS avatar, u.mobile AS mobile
		FROM
		    statistics s
		JOIN
		    user_info u ON s.uid = u.id
		WHERE
		    s.period = ?  -- 可以按需添加查询条件
		ORDER BY
		    s.accomplish_count DESC  -- 按完成数降序排列
		LIMIT ? OFFSET ?  -- 分页参数
	*/
	result := rr.db_.Table("statistics s").
		Select(`
            s.uid, s.accomplish_count,
            u.nickname AS username, u.avatar_url AS avatar, u.mobile AS mobile`).
		Joins("JOIN user_info u ON s.uid = u.id").
		Where("s.period = ?", period).
		Order("s.accomplish_count DESC").
		Limit(limit).
		Scan(&lb_datas)
	if result.Error != nil {
		logrus.Errorf("query failed, err:%v", result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "no data found")
	}

	lb_list := make([]*pb.LeaderboardUserInfo, 0, len(lb_datas))
	for _, lb_data := range lb_datas {
		lb_list = append(lb_list, &pb.LeaderboardUserInfo{
			Uid:      lb_data.Uid,
			Score:    lb_data.AccomplishCount,
			UserName: lb_data.Username,
			Avatar:   lb_data.Avatar,
			Mobile:   lb_data.Mobile,
		})
	}
	return lb_list, nil
}
func (rr *RecordRepo) QueryDailyAccomplishLeaderboard(limit int) ([]*pb.LeaderboardUserInfo, error) {
	var lb_datas []leaderboardData
	/*
		SELECT
		    us.uid, COUNT(DISTINCT us.problem_id) AS today_accomplish_count,
		    u.nickname AS username, u.avatar_url AS avatar, u.mobile AS mobile
		FROM
		    user_solution us
		JOIN
		    user_info u ON us.uid = u.id
		WHERE
		    DATE(us.create_at) = CURRENT_DATE() // 只统计今天的
		GROUP BY
		    us.uid, u.nickname, u.avatar_url, u.mobile
		ORDER BY
		    today_accomplish_count DESC
		LIMIT 200;
	*/
	result := rr.db_.Table("user_solution us").
		Select(`
            us.uid, COUNT(DISTINCT us.problem_id) AS today_accomplish_count,
            u.nickname AS username, u.avatar_url AS avatar, u.mobile AS mobile`).
		Joins("JOIN user_info u ON us.uid = u.id").
		Where("DATE(us.create_at) = CURRENT_DATE()").
		Group("us.uid, u.nickname, u.avatar_url, u.mobile").
		Order("today_accomplish_count DESC").
		Limit(limit).
		Scan(&lb_datas)
	if result.Error != nil {
		logrus.Errorf("query failed, err:%v", result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "no data found")
	}
	lb_list := make([]*pb.LeaderboardUserInfo, 0, len(lb_datas))
	for _, lb_data := range lb_datas {
		lb_list = append(lb_list, &pb.LeaderboardUserInfo{
			Uid:      lb_data.Uid,
			UserName: lb_data.Username,
			Avatar:   lb_data.Avatar,
			Score:    lb_data.AccomplishCount,
			Mobile:   lb_data.Mobile,
		})
	}
	return lb_list, nil
}
func (rr *RecordRepo) SynchronizeLeaderboard(lb_list []*pb.LeaderboardUserInfo, leaderboardKey string, leaderboardKeyTTL time.Duration) error {
	// 使用Pipeline批量操作
	ctx := context.Background()
	pipe := rr.rdb_.Pipeline()

	// 删除旧榜
	pipe.Del(ctx, leaderboardKey)
	pipe.Del(ctx, global.LeaderboardUserInfoKey)

	// 批量添加新数据
	for _, item := range lb_list {
		pipe.ZAdd(ctx, leaderboardKey, redis.Z{
			Score:  float64(item.Score),
			Member: item.Uid,
		})
		userInfo, err := protojson.Marshal(item)
		if err != nil {
			logrus.Errorf("failed to marshal user info: %v", err)
			return err
		}
		pipe.HSet(ctx, global.LeaderboardUserInfoKey, strconv.FormatInt(item.Uid, 10), userInfo)
	}
	pipe.Expire(ctx, leaderboardKey, leaderboardKeyTTL)

	// 执行Pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		logrus.Errorf("failed to synchronize leaderboard: %v", err)
		return fmt.Errorf("failed to synchronize leaderboard: %w", err)
	}
	return nil
}
