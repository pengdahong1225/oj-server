package data

import (
	"context"
	"fmt"
	"io"
	"log"
	"oj-server/common/consts"
	"oj-server/svr/problem/internal/config"
	"oj-server/svr/problem/internal/model"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type NoticeRepo struct {
	db_  *gorm.DB
	rdb_ *redis.Client
}

func NewNoticeRepo() (*NoticeRepo, error) {
	// mysql
	mysql_cfg := config.MyConfig.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysql_cfg.User,
		mysql_cfg.Pwd, mysql_cfg.Host, mysql_cfg.Port, mysql_cfg.Db)
	timer := time.Now().Format("2006_01_02")
	filePath := fmt.Sprintf("%s/orm.%s.log", consts.LogPath, timer)
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
	redis_cfg := config.MyConfig.Redis
	dsn = fmt.Sprintf("%s:%d", redis_cfg.Host, redis_cfg.Port)
	rdb_ := redis.NewClient(&redis.Options{
		Addr:    dsn,
		Network: "tcp",
	})
	st := rdb_.Ping(context.Background())
	if st.Err() != nil {
		return nil, st.Err()
	}

	return &NoticeRepo{
		db_:  db_,
		rdb_: rdb_,
	}, nil
}

// CreateNotice 新建公告
func (nr *NoticeRepo) CreateNotice(notice *model.Notice) error {
	result := nr.db_.Create(notice)
	if result.Error != nil {
		logrus.Errorf("create notice failed: %s", result.Error.Error())
		return status.Errorf(codes.Internal, "create notice failed")
	}
	return nil
}

// DeleteNotice 删除公告
func (nr *NoticeRepo) DeleteNotice(id int64) error {
	result := nr.db_.Where("id = ?", id).Delete(&model.Notice{})
	if result.Error != nil {
		logrus.Errorf("delete notice failed: %s", result.Error.Error())
		return status.Errorf(codes.Internal, "delete notice failed")
	}
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "notice not found")
	}
	return nil
}

// UpdateNotice 更新公告标题和内容
func (nr *NoticeRepo) UpdateNotice(notice *model.Notice) error {
	result := nr.db_.Model(&model.Notice{}).Where("id = ?", notice.ID).
		Updates(map[string]any{
			"title":   notice.Title,
			"content": notice.Content,
		})
	if result.Error != nil {
		logrus.Errorf("update notice failed: %s", result.Error.Error())
		return status.Errorf(codes.Internal, "update notice failed")
	}
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "notice not found")
	}
	return nil
}

// UpdateNoticeStatus 更新公告状态（1 发布 / 0 隐藏）
func (nr *NoticeRepo) UpdateNoticeStatus(id int64, st int32) error {
	result := nr.db_.Model(&model.Notice{}).Where("id = ?", id).Update("status", st)
	if result.Error != nil {
		logrus.Errorf("update notice status failed: %s", result.Error.Error())
		return status.Errorf(codes.Internal, "update notice status failed")
	}
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "notice not found")
	}
	return nil
}

// QueryNoticeList 分页查询公告列表
func (nr *NoticeRepo) QueryNoticeList(page, pageSize int) (int64, []model.Notice, error) {
	var count int64
	result := nr.db_.Model(&model.Notice{}).Count(&count)
	if result.Error != nil {
		logrus.Errorf("query notice count failed: %s", result.Error.Error())
		return -1, nil, status.Errorf(codes.Internal, "query notice count failed")
	}
	if count == 0 {
		return 0, nil, nil
	}

	offset := (page - 1) * pageSize
	var notices []model.Notice
	result = nr.db_.Model(&model.Notice{}).
		Select("id, title, create_at, status, create_by").
		Order("id desc").
		Offset(offset).
		Limit(pageSize).
		Find(&notices)
	if result.Error != nil {
		logrus.Errorf("query notice list failed: %s", result.Error.Error())
		return -1, nil, status.Errorf(codes.Internal, "query notice list failed")
	}
	return count, notices, nil
}

// QueryNotice 查询单个公告详情
func (nr *NoticeRepo) QueryNotice(id int64) (*model.Notice, error) {
	var notice model.Notice
	result := nr.db_.Where("id = ?", id).Find(&notice)
	if result.Error != nil {
		logrus.Errorf("query notice failed: %s", result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query notice failed")
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "notice not found")
	}
	return &notice, nil
}
