package domain

import (
	"oj-server/module/model"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 查询用户的历史提交记录
func (md *MysqlDB) QuerySubmitRecordList(uid int64, pageSize, offset int) (int64, []model.SubmitRecord, error) {
	/*
		select id, created_at, problem_name, status, lang from user_submit_record
		where uid = ?
		order by created_at desc
		offset off_set
		limit page_size;
	*/
	var count int64 = 0
	result := md.db_.Model(&model.SubmitRecord{}).Where("uid = ?", uid).Count(&count)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return 0, nil, status.Errorf(codes.Internal, "查询提交记录失败")
	}
	var records []model.SubmitRecord
	result = md.db_.Where("uid = ?", uid).Order("created_at desc").Offset(offset).Limit(pageSize).Find(&records)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return 0, nil, status.Errorf(codes.Internal, "查询提交记录失败")
	}
	return count, records, nil
}

// 查询某项提交记录的详细信息
func (md *MysqlDB) QuerySubmitRecord(id int64) (*model.SubmitRecord, error) {
	var record model.SubmitRecord
	result := md.db_.Where("id = ?", id).First(&record)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query failed")
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "record not found")
	}
	return &record, nil
}
