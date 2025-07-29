package domain

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"oj-server/module/model"
)

func (md *MysqlDB) CreateProblem(problem *model.Problem) (int64, error) {
	result := md.db_.Create(problem)
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
func (md *MysqlDB) QueryProblemList(page, pageSize int, keyword, tag string) (int64, []model.Problem, error) {
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
		result = md.db_.Model(&model.Problem{}).Where("title LIKE ?", name).Count(&count)
	} else {
		result = md.db_.Model(&model.Problem{}).Where("title LIKE ?", name).Where(query).Count(&count)
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
		result = md.db_.Select("id,title,level,tags,create_at,create_by").Where("title LIKE ?", name).Order("id").Offset(offSet).Limit(pageSize).Find(&problemList)
	} else {
		result = md.db_.Select("id,title,level,tags,create_at,create_by").Where("title LIKE ?", name).Where(query).Order("id").Offset(offSet).Limit(pageSize).Find(&problemList)
	}
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return -1, nil, status.Errorf(codes.Internal, "query failed")
	}
	return count, problemList, nil
}

func (md *MysqlDB) QueryProblemData(id int64) (*model.Problem, error) {
	var problem model.Problem
	result := md.db_.Where("id=?", id).Find(&problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query failed")
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "problem not found")
	}
	return &problem, nil
}

func (md *MysqlDB) UpdateProblem(problem *model.Problem) error {
	result := md.db_.Model(&model.Problem{}).Where("id=?", problem.ID).Updates(problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return status.Errorf(codes.Internal, "query failed")
	}
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "problem not found")
	}
	return nil
}

func (md *MysqlDB) UpdateProblemStatus(id int64, st int32) error {
	result := md.db_.Model(&model.Problem{}).Where("id=?", id).Update("status", st)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return status.Errorf(codes.Internal, "query failed")
	}
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "query failed")
	}
	return nil
}

func (md *MysqlDB) DeleteProblem(id int64) error {
	result := md.db_.Where("id=?", id).Delete(&model.Problem{})
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return status.Errorf(codes.Internal, "delete problem failed")
	}
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "problem not found")
	}
	return nil
}
