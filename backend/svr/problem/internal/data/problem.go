package data

import (
	"context"
	"encoding/json"
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
	"oj-server/svr/problem/internal/configs"
	"oj-server/svr/problem/internal/model"
	"os"
	"time"
)

type ProblemRepo struct {
	db_  *gorm.DB
	rdb_ *redis.Client
}

func NewProblemRepo() (*ProblemRepo, error) {
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

	return &ProblemRepo{
		db_:  db_,
		rdb_: rdb_,
	}, nil
}

func (pr *ProblemRepo) CreateProblem(problem *model.Problem) (int64, error) {
	result := pr.db_.Create(problem)
	if result.Error != nil {
		logrus.Errorf("create problem failed: %s", result.Error.Error())
		return -1, status.Errorf(codes.Internal, "create problem failed")
	}
	return problem.ID, nil
}

// 分页查询题库列表
// 查询{id，title，level，tags}
// @page 页码
// @page_size 单页数量
// @keyword 关键字
// @tag 标签
func (pr *ProblemRepo) QueryProblemList(page, pageSize int, keyword, tag string) (int64, []model.Problem, error) {
	/**
	select COUNT(*) AS count
	from problem
	where title like '%name%' AND JSON_CONTAINS(tags, '"哈希表"');
	*/
	var count int64 = 0
	query := pr.db_.Model(&model.Problem{})
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if tag != "" {
		str := fmt.Sprintf(`JSON_CONTAINS(tags, '"%s"')`, tag)
		query = query.Where(str)
	}
	result := query.Count(&count)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return -1, nil, status.Errorf(codes.Internal, "query failed")
	}

	/*
		select id,title,level,tags,create_at,create_by from problem
		where title like '%name%' AND JSON_CONTAINS(tags, '"哈希表')
		order by id
		offset off_set
		limit page_size;
	*/
	offSet := (page - 1) * pageSize
	var problemList []model.Problem
	query = pr.db_.Model(&model.Problem{}).Select("id,title,level,tags,create_at,update_at,status")
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if tag != "" {
		query = query.Where("JSON_CONTAINS(tags, ?)", tag)
	}
	result = query.Order("id").Offset(offSet).Limit(pageSize).Find(&problemList)
	if result.Error != nil {
		logrus.Errorf("query problem list failed: %s", result.Error.Error())
		return -1, nil, status.Errorf(codes.Internal, "query failed")
	}
	return count, problemList, nil
}

func (pr *ProblemRepo) QueryProblemData(id int64) (*model.Problem, error) {
	var problem model.Problem
	result := pr.db_.Where("id=?", id).Find(&problem)
	if result.Error != nil {
		logrus.Errorf("query problem failed: %s", result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query failed")
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "problem not found")
	}
	return &problem, nil
}

func (pr *ProblemRepo) UpdateProblem(problem *model.Problem) error {
	result := pr.db_.Model(&model.Problem{}).Where("id=?", problem.ID).Updates(problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return status.Errorf(codes.Internal, "query failed")
	}
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "problem not found")
	}
	return nil
}

func (pr *ProblemRepo) UpdateProblemStatus(id int64, st int32) error {
	result := pr.db_.Model(&model.Problem{}).Where("id=?", id).Update("status", st)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return status.Errorf(codes.Internal, "query failed")
	}
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "query failed")
	}
	return nil
}

func (pr *ProblemRepo) DeleteProblem(id int64) error {
	//result := pr.db_.Where("id=?", id).Update("delete_at", time.Now().String()) // 软删除
	result := pr.db_.Where("id=?", id).Delete(&model.Problem{})
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return status.Errorf(codes.Internal, "delete problem failed")
	}
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "problem not found")
	}
	return nil
}

// 查询标签列表
func (pr *ProblemRepo) QueryTagList() ([]string, error) {
	ctx := context.Background()

	// 1. 先尝试从缓存读取
	tagList, err := pr.rdb_.SMembers(ctx, global.TagListKey).Result()
	if err == nil && len(tagList) > 0 {
		logrus.Debugf("tag list from cache: %v", tagList)
		return tagList, nil
	}

	logrus.Warn("tag list cache miss, querying database ...")

	// 2. 缓存没有 → 查询数据库的所有 tags 字段
	var tagsJsonList []string
	result := pr.db_.Model(&model.Problem{}).Pluck("tags", &tagsJsonList)
	if result.Error != nil {
		logrus.Errorf("query tag list failed: %s", result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query failed")
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	// 3.用 map 做去重, 解析 JSON 数组
	tagSet := make(map[string]struct{})
	for _, jsonStr := range tagsJsonList {
		if jsonStr == "" {
			continue
		}
		var arr []string
		if err = json.Unmarshal([]byte(jsonStr), &arr); err != nil {
			logrus.Errorf("failed to unmarshal tags: %s", err.Error())
			continue
		}
		for _, t := range arr {
			tagSet[t] = struct{}{}
		}
	}
	tagList = make([]string, 0, len(tagSet))
	for t := range tagSet {
		tagList = append(tagList, t)
	}
	// 4. 持久化到 Redis（Set 类型）
	if len(tagList) > 0 {
		members := make([]interface{}, 0, len(tagList))
		for _, t := range tagList {
			members = append(members, t)
		}
		pr.rdb_.SAdd(ctx, global.TagListKey, members...)
	}

	logrus.Debugf("tag list final: %v", tagList)
	return tagList, nil
}

func (pr *ProblemRepo) Lock(key string, ttl time.Duration) (bool, error) {
	return pr.rdb_.SetNX(context.Background(), key, "locked", ttl).Result()
}
func (pr *ProblemRepo) UnLock(key string) error {
	return pr.rdb_.Del(context.Background(), key).Err()
}
