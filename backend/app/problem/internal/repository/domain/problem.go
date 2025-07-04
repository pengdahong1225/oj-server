package domain

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"oj-server/app/common/errs"
	"oj-server/app/problem/internal/repository/model"
)

func (md *MysqlDB) CreateProblem(problem *model.Problem) (int64, error) {
	return -1, nil
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
		return -1, nil, errs.QueryFailed
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
		return -1, nil, errs.QueryFailed
	}
	return count, problemList, nil
}

func (md *MysqlDB) QueryProblemData(id int64) (*model.Problem, error) {
	var problem model.Problem
	result := md.db_.Where("id=?", id).Find(&problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, errs.QueryFailed
	}
	if result.RowsAffected == 0 {
		return nil, errs.NotFound
	}
	return &problem, nil
}
