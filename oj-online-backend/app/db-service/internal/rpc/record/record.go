package record

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/internal/rpc"
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/internal/svc/mysql"
	"github.com/pengdahong1225/Oj-Online-Server/proto/pb"
	"github.com/sirupsen/logrus"
)

type RecordServer struct {
	pb.UnimplementedRecordServiceServer
}

func (receiver *RecordServer) SaveUserSubmitRecord(ctx context.Context, request *pb.SaveUserSubmitRecordRequest) (*empty.Empty, error) {
	/*
		insert into user_submit_record_xxx
		(uid, problem_id, code, result, lang)
		values
		(?, ?, ?, ?, ?);
	*/
	db := mysql.Instance()
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
			return nil, rpc.InsertFailed
		}
	}
	result := db.Table(record.TableName(request.Stamp)).Create(record)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc.InsertFailed
	}

	return &empty.Empty{}, nil
}

func (receiver *RecordServer) GetUserSubmitRecord(ctx context.Context, request *pb.GetUserSubmitRecordRequest) (*pb.GetUserSubmitRecordResponse, error) {
	/*
		select * from user_submit_record_xx
		where uid = ? and stamp = ?;
	*/
	db := mysql.Instance()
	r := &mysql.SubMitRecord{}
	if !db.Migrator().HasTable(r.TableName(request.Stamp)) {
		return nil, rpc.NotFound
	}

	var records []mysql.SubMitRecord
	result := db.Where("uid = ?", request.UserId).Find(&records)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc.QueryFailed
	}
	if result.RowsAffected == 0 {
		return nil, rpc.NotFound
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
