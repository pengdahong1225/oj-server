package record

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/svc/mysql"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
)

type RecordServer struct {
	pb.UnimplementedRecordServiceServer
}

// SaveUserSubmitRecord
// 1.更新用户提交记录表
// 2.更新(用户,题目)AC表
// 3.更新用户解题情况统计表
func (receiver *RecordServer) SaveUserSubmitRecord(ctx context.Context, request *pb.SaveUserSubmitRecordRequest) (*empty.Empty, error) {
	tx := mysql.Instance().Begin()
	if tx.Error != nil {
		logrus.Errorln(tx.Error.Error())
		return nil, tx.Error
	}

	// todo 更新用户提交记录表
	/*
		insert into user_submit_record
		(uid, problem_id, problem_name, status, code, result, lang)
		values
		(?, ?, ?, ?, ?, ?, ?);
	*/
	record := &mysql.SubmitRecord{
		Uid:         request.Data.Uid,
		ProblemID:   request.Data.ProblemId,
		ProblemName: request.Data.ProblemName,
		Status:      request.Data.Status,
		Code:        request.Data.Code,
		Result:      request.Data.Result,
		Lang:        request.Data.Lang,
	}
	result := tx.Create(record)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		tx.Rollback()
		return nil, rpc_api.InsertFailed
	}

	// todo 更新(用户,题目)AC表
	/*
		insert into user_solution(uid,problem_id)
		values(?,?)
	*/
	var repeatedAc = true
	if request.Data.Status == "Accepted" {
		data := mysql.UserSolution{}
		result = tx.Where("uid=? and problem_id=?", request.Data.Uid, request.Data.ProblemId).Find(&data)
		if result.RowsAffected == 0 {
			repeatedAc = false

			data.Uid = request.Data.Uid
			data.ProblemID = request.Data.ProblemId

			result = tx.Create(&data)
			if result.Error != nil {
				logrus.Errorln(result.Error.Error())
				tx.Rollback()
				return nil, rpc_api.InsertFailed
			}
		}
	}

	// todo 更新用户解题情况统计表
	data := mysql.Statistics{
		Uid: request.Data.Uid,
	}
	result = tx.FirstOrCreate(&data)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc_api.QueryFailed
	}
	data.SubmitCount += 1
	if !repeatedAc && request.Data.Status == "Accepted" {
		data.AccomplishCount += 1
		switch request.Data.ProblemLevel {
		case 1:
			data.EasyProblemCount += 1
		case 2:
			data.MediumProblemCount += 1
		case 3:
			data.HardProblemCount += 1
		}
	}
	result = tx.Where("uid=?", data.Uid).Save(&data)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		tx.Rollback()
		return nil, rpc_api.UpdateFailed
	}

	tx.Commit()
	return &empty.Empty{}, nil
}

// GetUserRecordList
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
	rsp := &pb.GetUserRecordListResponse{}

	var count int64 = 0
	result := db.Model(&mysql.SubmitRecord{}).Where("uid = ?", request.Uid).Count(&count)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc_api.QueryFailed
	}
	rsp.Total = int32(count)
	if count == 0 {
		return rsp, nil
	}

	var records []mysql.SubmitRecord
	result = db.Where("uid = ?", request.Uid).Order("created_at desc").Offset(offSet).Limit(int(request.PageSize)).Find(&records)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc_api.QueryFailed
	}
	if result.RowsAffected == 0 {
		return nil, rpc_api.NotFound
	}

	for _, record := range records {
		rsp.Data = append(rsp.Data, &pb.UserSubmitRecord{
			Id:          int64(record.ID),
			CreatedAt:   record.CreatedAt.Unix(),
			ProblemName: record.ProblemName,
			Status:      record.Status,
			Lang:        record.Lang,
		})
	}

	return rsp, nil
}

func (receiver *RecordServer) GetUserRecord(ctx context.Context, request *pb.GetUserRecordRequest) (*pb.GetUserRecordResponse, error) {
	db := mysql.Instance()
	var record mysql.SubmitRecord
	result := db.Where("id = ?", request.Id).First(&record)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc_api.QueryFailed
	}
	if result.RowsAffected == 0 {
		return nil, rpc_api.NotFound
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
