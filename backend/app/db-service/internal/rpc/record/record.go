package record

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/rpc"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/svc/mysql"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
)

type RecordServer struct {
	pb.UnimplementedRecordServiceServer
}

func (receiver *RecordServer) SaveUserSubmitRecord(ctx context.Context, request *pb.SaveUserSubmitRecordRequest) (*empty.Empty, error) {
	/*
		insert into user_submit_record
		(uid, problem_id, problem_name, status, code, result, lang)
		values
		(?, ?, ?, ?, ?, ?, ?);
	*/
	db := mysql.Instance()
	record := &mysql.SubmitRecord{
		Uid:         request.Data.Uid,
		ProblemID:   request.Data.ProblemId,
		ProblemName: request.Data.ProblemName,
		Status:      request.Data.Status,
		Code:        request.Data.Code,
		Result:      request.Data.Result,
		Lang:        request.Data.Lang,
	}
	//if !db.Migrator().HasTable(record.TableName(request.Stamp)) {
	//	err := db.Table(record.TableName(request.Stamp)).AutoMigrate(record)
	//	if err != nil {
	//		logrus.Errorln(err.Error())
	//		return nil, rpc.InsertFailed
	//	}
	//}
	result := db.Create(record)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc.InsertFailed
	}

	return &empty.Empty{}, nil
}

// GetUserSubmitRecord
// 分页查询用户的提交记录
// @uid
// @page
// @pageSize
func (receiver *RecordServer) GetUserRecordList(ctx context.Context, request *pb.GetUserRecordListRequest) (*pb.GetUserRecordListResponse, error) {
	/*
		select id, created_at, problem_name, status, lang from user_submit_record
		where uid = ?
		order by created_at desc
		offset off_set
		limit page_size;
	*/
	db := mysql.Instance()
	offSet := int((request.Page - 1) * request.PageSize)

	var records []mysql.SubmitRecord
	result := db.Where("uid = ?", request.Uid).Order("created_at desc").Offset(offSet).Limit(int(request.PageSize)).Find(&records)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc.QueryFailed
	}
	if result.RowsAffected == 0 {
		return nil, rpc.NotFound
	}

	list := make([]*pb.UserSubmitRecord, 0, len(records))
	for _, record := range records {
		list = append(list, &pb.UserSubmitRecord{
			Id:          int64(record.ID),
			CreatedAt:   record.CreatedAt.Unix(),
			ProblemName: record.ProblemName,
			Status:      record.Status,
			Lang:        record.Lang,
		})
	}

	return &pb.GetUserRecordListResponse{Data: list}, nil
}

func (receiver *RecordServer) GetUserRecord(ctx context.Context, request *pb.GetUserRecordRequest) (*pb.GetUserRecordResponse, error) {
	db := mysql.Instance()
	var record mysql.SubmitRecord
	result := db.Where("id = ?", request.Id).First(&record)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc.QueryFailed
	}
	if result.RowsAffected == 0 {
		return nil, rpc.NotFound
	}
	return &pb.GetUserRecordResponse{Data: &pb.UserSubmitRecord{
		Id:        int64(record.ID),
		CreatedAt: record.CreatedAt.Unix(),

		ProblemId:   record.ProblemID,
		ProblemName: record.ProblemName,
		Status:      record.Status,
		Lang:        record.Lang,
		Code:        record.Code,
		Result:      record.Result,
	}}, nil
}
