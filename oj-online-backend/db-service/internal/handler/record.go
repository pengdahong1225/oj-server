package handler

import (
	"context"
	pb "db-service/internal/proto"
	mysql "db-service/services/mysql"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
)

func (receiver *DBServiceServer) CreateUserSubmitRecord(ctx context.Context, request *pb.CreateUserSubmitRecordRequest) (*empty.Empty, error) {
	/*
		insert into user_submit_record_xxx
		(uid, problem_id, code, result, lang)
		values
		(?, ?, ?, ?, ?);
	*/
	db := mysql.DB
	record := &mysql.SubMitRecord{
		Uid:       request.UserId,
		ProblemID: request.ProblemId,
		Code:      request.Code,
		Result:    request.Result,
		Lang:      request.Lang,
	}
	if !db.Migrator().HasTable(record.TableName(request.Stamp)) {
		err := db.Table(record.TableName(request.Stamp)).AutoMigrate(record)
		if err != nil {
			logrus.Errorln(err.Error())
			return nil, InsertFailed
		}
	}
	result := db.Table(record.TableName(request.Stamp)).Create(record)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, InsertFailed
	}

	return &empty.Empty{}, nil
}

func (receiver *DBServiceServer) GetUserSubmitRecord(ctx context.Context, request *pb.GetUserSubmitRecordRequest) (*pb.GetUserSubmitRecordResponse, error) {
	/*
		select * from user_submit_record_xx
		where uid = ? and stamp = ?;
	*/
	db := mysql.DB
	r := &mysql.SubMitRecord{}
	if !db.Migrator().HasTable(r.TableName(request.Stamp)) {
		return nil, NotFound
	}

	var records []mysql.SubMitRecord
	result := db.Where("uid = ?", request.UserId).Find(&records)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, QueryField
	}
	if result.RowsAffected == 0 {
		return nil, NotFound
	}

	data := make([]*pb.UserSubmitRecord, 0, len(records))
	for _, record := range records {
		data = append(data, &pb.UserSubmitRecord{
			Uid:       record.Uid,
			ProblemId: record.ProblemID,
			Code:      record.Code,
			Result:    record.Result,
			Lang:      record.Lang,
			Stamp:     record.CreatedAt.Unix(),
		})
	}

	return &pb.GetUserSubmitRecordResponse{Data: data}, nil
}
