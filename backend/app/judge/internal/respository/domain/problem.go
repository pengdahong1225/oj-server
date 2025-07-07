package domain

import (
	"oj-server/module/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/sirupsen/logrus"
)

func (r *MysqlDB) QueryProblemData(id int64) (*model.Problem, error) {
	var problem model.Problem
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
