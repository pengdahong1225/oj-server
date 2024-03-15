package handler

import (
	"context"
	"db-service/global"
	"db-service/internal/models"
	pb "db-service/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (receiver *DBServiceServer) GetUserSubmitRecord(ctx context.Context, request *pb.GetUserSubmitRecordRequest) (*pb.GetUserSubmitRecordResponse, error) {
	// 查询用户提交记录
	// select * form user_submit where user_id = id
	var recordList []models.UserSubMit
	result := global.DBInstance.Where("id=?", request.UserId).Find(&recordList)
	if result.Error != nil {
		logrus.Errorln("GetUserSubmitRecord error", result.Error)
		return nil, status.Errorf(codes.Internal, "GetUserSubmitRecord error")
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "GetUserSubmitRecord not found")
	}

	list := make([]*pb.SubmitRecord, len(recordList))
	for i, record := range recordList {
		recordList[i].UserID = record.UserID
		recordList[i].QuestionID = record.QuestionID
		recordList[i].Code = record.Code
		recordList[i].Result = record.Result
		recordList[i].Lang = record.Lang
		recordList[i].CreateAt = record.CreateAt
	}
	rsp := &pb.GetUserSubmitRecordResponse{
		Data: list,
	}
	return rsp, nil
}

func (receiver *DBServiceServer) UpdateUserSubmitRecord(ctx context.Context, request *pb.UpdateUserSubmitRecordRequest) (*empty.Empty, error) {
	record := models.UserSubMit{
		UserID:     request.UserId,
		QuestionID: request.QuestionId,
		Code:       request.Code,
		Result:     request.Result,
		Lang:       request.Lang,
	}
	result := global.DBInstance.Create(&record)
	if result.Error != nil {
		logrus.Errorf("插入记录失败：%s", result.Error.Error())
		return nil, status.Errorf(codes.Internal, "插入记录失败：%s", result.Error.Error())
	} else if result.RowsAffected == 0 {
		logrus.Errorln("插入记录失败:result.RowsAffected == 0")
		return nil, status.Errorf(codes.NotFound, "插入记录失败")
	}
	return &empty.Empty{}, nil
}
