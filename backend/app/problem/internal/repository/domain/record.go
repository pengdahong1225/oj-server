package domain

import (
	"oj-server/module/model"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

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
func (md *MysqlDB) QueryUserRecord(uid int64, offset int, pageSize int) (int64, []model.SubmitRecord, error) {
}
