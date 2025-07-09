package service

import (
	"oj-server/proto/pb"
	"context"
	"github.com/sirupsen/logrus"
	"oj-server/module/model"
)

// GetSubmitRecordList
// 分页查询用户的提交记录
// @uid
// @page
// @pageSize
func (ps *ProblemService) GetSubmitRecordList(ctx context.Context, in *pb.GetSubmitRecordListRequest) (*pb.GetSubmitRecordListResponse, error) {
	offSet := int((in.Page - 1) * in.PageSize)
	count, records, err := ps.db.QuerySubmitRecordList(in.Uid, int(in.PageSize), offSet)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, err
	}

	resp := &pb.GetSubmitRecordListResponse{
		Total: int32(count),
		Data:  make([]*pb.SubmitRecord, len(records)),
	}
	for i, record := range records {
		resp.Data[i] = &pb.SubmitRecord{
			Id:          int64(record.ID),
			CreatedAt:   record.CreatedAt.Unix(),
			ProblemName: record.ProblemName,
			Lang:        record.Lang,
			Status:      record.Status,
		}
	}

	return resp, nil
}

func (ps *ProblemService) GetSubmitRecordData(ctx context.Context, in *pb.GetSubmitRecordRequest) (*pb.GetSubmitRecordResponse, error) {
	var record model.SubmitRecord
	result := db.Where("id = ?", request.Id).First(&record)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, errs.QueryFailed
	}
	if result.RowsAffected == 0 {
		return nil, errs.NotFound
	}
	return &pb.GetUserRecordResponse{Data: &pb.UserSubmitRecord{
		Id:        int64(record.ID),
		CreatedAt: record.CreatedAt.Unix(),

		Uid:         record.Uid,
		UserName:    record.UserName,
		ProblemId:   record.ProblemID,
		ProblemName: record.ProblemName,
		Status:      record.Status,
		Lang:        record.Lang,
		Code:        record.Code,
		Result:      record.Result,
	}}, nil
}
