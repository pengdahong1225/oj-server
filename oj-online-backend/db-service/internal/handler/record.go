package handler

import (
	"context"
	"db-service/internal/models"
	pb2 "db-service/internal/proto"
	"db-service/services/dao/mysql"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (receiver *DBServiceServer) GetUserSubmitRecord(ctx context.Context, request *pb2.GetUserSubmitRecordRequest) (*pb2.GetUserSubmitRecordResponse, error) {
	db := mysql.DB
	// 查询用户提交记录
	// select * form user_submit where user_id = id
	var recordList []models.UserSubMit
	result := db.Where("id=?", request.UserId).Find(&recordList)
	if result.Error != nil {
		logrus.Errorln("GetUserSubmitRecord error", result.Error)
		return nil, status.Errorf(codes.Internal, "GetUserSubmitRecord error")
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "GetUserSubmitRecord not found")
	}

	list := make([]*pb2.SubmitRecord, len(recordList))
	for i, record := range recordList {
		recordList[i].UserID = record.UserID
		recordList[i].QuestionID = record.QuestionID
		recordList[i].Code = record.Code
		recordList[i].Result = record.Result
		recordList[i].Lang = record.Lang
		recordList[i].CreateAt = record.CreateAt
	}
	rsp := &pb2.GetUserSubmitRecordResponse{
		Data: list,
	}
	return rsp, nil
}

func (receiver *DBServiceServer) UpdateUserSubmitRecord(ctx context.Context, request *pb2.UpdateUserSubmitRecordRequest) (*empty.Empty, error) {
	db := mysql.DB
	record := models.UserSubMit{
		UserID:     request.UserId,
		QuestionID: request.QuestionId,
		Code:       request.Code,
		Result:     request.Result,
		Lang:       request.Lang,
	}
	result := db.Create(&record)
	if result.Error != nil {
		logrus.Errorf("插入记录失败：%s", result.Error.Error())
		return nil, status.Errorf(codes.Internal, "插入记录失败：%s", result.Error.Error())
	} else if result.RowsAffected == 0 {
		logrus.Errorln("插入记录失败:result.RowsAffected == 0")
		return nil, status.Errorf(codes.NotFound, "插入记录失败")
	}
	return &empty.Empty{}, nil
}