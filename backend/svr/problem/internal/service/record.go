package service

import (
	"context"

	"github.com/sirupsen/logrus"
	"oj-server/module/proto/pb"
)

// 分页查询用户的提交记录
// @uid
// @page
// @pageSize
func (ps *ProblemService) GetSubmitRecordList(ctx context.Context, in *pb.GetSubmitRecordListRequest) (*pb.GetSubmitRecordListResponse, error) {
	offSet := int((in.Page - 1) * in.PageSize)
	count, records, err := ps.rc.QuerySubmitRecordList(in.Uid, int(in.PageSize), offSet)
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
	record, err := ps.rc.QuerySubmitRecord(in.Id)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, err
	}
	resp := &pb.GetSubmitRecordResponse{
		Data: record.Transform(),
	}

	return resp, nil
}
