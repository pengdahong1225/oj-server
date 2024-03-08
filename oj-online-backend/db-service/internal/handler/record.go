package handler

import (
	"context"
	"db-service/global"
	"db-service/internal/models"
	pb "db-service/proto"
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
