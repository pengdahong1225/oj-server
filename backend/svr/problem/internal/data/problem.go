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
)

type ProblemRepo struct {
	db_  *gorm.DB
	rdb_ *redis.Client
}

func NewProblemRepo() (*ProblemRepo, error) {
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

	return &ProblemRepo{
		db_:  db_,
		rdb_: rdb_,
	}, nil
}

func (pr *ProblemRepo) CreateProblem(problem *model.Problem) (int64, error) {
	result := pr.db_.Create(problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return -1, status.Errorf(codes.Internal, "create problem failed")
	}
	return problem.ID, nil
}

// QueryProblemList 分页查询题库列表
// 查询{id，title，level，tags}
// @page 页码
// @page_size 单页数量
// @keyword 关键字
// @tag 标签
func (pr *ProblemRepo) QueryProblemList(page, pageSize int, keyword, tag string) (int64, []model.Problem, error) {
	name := "%" + keyword + "%"
	offSet := (page - 1) * pageSize
	query := fmt.Sprintf(`JSON_CONTAINS(tags, '"%s"')`, tag)

	logrus.Debugf("query conditions: %s\n", query)
	/*
		select COUNT(*) AS count from problem
		where title like '%name%' AND JSON_CONTAINS(tags, '"哈希表"');
	*/
	var result *gorm.DB
	var count int64 = 0
	if tag == "" {
		result = pr.db_.Model(&model.Problem{}).Where("title LIKE ?", name).Count(&count)
	} else {
		result = pr.db_.Model(&model.Problem{}).Where("title LIKE ?", name).Where(query).Count(&count)
	}
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
	var problemList []model.Problem
	if tag == "" {
		result = pr.db_.Select("id,title,level,tags,create_at,create_by").Where("title LIKE ?", name).Order("id").Offset(offSet).Limit(pageSize).Find(&problemList)
	} else {
		result = pr.db_.Select("id,title,level,tags,create_at,create_by").Where("title LIKE ?", name).Where(query).Order("id").Offset(offSet).Limit(pageSize).Find(&problemList)
	}
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return -1, nil, status.Errorf(codes.Internal, "query failed")
	}
	return count, problemList, nil
}

func (pr *ProblemRepo) QueryProblemData(id int64) (*model.Problem, error) {
	var problem model.Problem
	result := pr.db_.Where("id=?", id).Find(&problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
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
// 要区分not found 和 err
func (pr *ProblemRepo) QueryTagList() ([]string, error) {
	var (
		tagList []string
		err     error
	)
	// 先查询缓存
	tagList, err = pr.rdb_.SMembers(context.Background(), global.TagListKey).Result()
	if err != nil {
		// 缓存中没有or查询缓存失败，则从数据库中查询
		/*
			select tags from problem
		*/
		result := pr.db_.Model(&model.Problem{}).Pluck("tags", &tagList)
		if result.Error != nil {
			logrus.Errorf("query tag list failed: %s", result.Error.Error())
			return nil, status.Errorf(codes.Internal, "query failed")
		}
		if result.RowsAffected == 0 {
			return nil, nil // not found
		}
		for _, tag := range tagList {
			// 添加到缓存
			pr.rdb_.SAdd(context.Background(), global.TagListKey, tag)
		}
	}

	return tagList, nil
}
