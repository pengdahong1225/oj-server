package handler

import (
	"context"
	"db-service/global"
	"db-service/internal/models"
	pb "db-service/proto"
	"db-service/utils"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// 题库列表，每次修改都更新同步到redis

func (receiver *DBServiceServer) GetQuestionData(ctx context.Context, request *pb.GetQuestionRequest) (*pb.GetQuestionResponse, error) {
	var question models.Question
	result := global.DBInstance.Where("id = ?", request.Id).Find(&question)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "GetQuestionData error: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "GetQuestionData[%d] error: %v", request.Id, "not found")
	}

	data := &pb.Question{
		Id:          question.ID,
		CreateAt:    timestamppb.New(question.CreateAt),
		Title:       question.Title,
		Description: question.Description,
		Level:       question.Level,
		Tags:        utils.SplitStringWithX(question.Tags, "#"),
	}
	return &pb.GetQuestionResponse{
		Data: data,
	}, nil
}

func (receiver *DBServiceServer) CreateQuestionData(ctx context.Context, request *pb.CreateQuestionRequest) (*pb.CreateQuestionResponse, error) {
	question := &models.Question{
		Title:       request.Data.Title,
		Description: request.Data.Description,
		Level:       request.Data.Level,
		Tags:        utils.SpliceStringWithX(request.Data.Tags, "#"),
	}
	result := global.DBInstance.Where("id = ?", question.ID)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "GetQuestionData error: %v", result.Error)
	}
	if result.RowsAffected > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "method CreateQuestionData error: %v", "already exists")
	}

	result = global.DBInstance.Create(question)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "method CreateQuestionData error: %v", result.Error)
	}

	// 插入成功 -- 将测试案例插入到redis
	updateTestCases(request.Data.Id, request.Data.TestCase)

	return &pb.CreateQuestionResponse{
		Id: question.ID,
	}, nil
}

func (receiver *DBServiceServer) UpdateQuestionData(ctx context.Context, request *pb.UpdateQuestionRequest) (*empty.Empty, error) {
	question := &models.Question{
		Title:       request.Data.Title,
		Description: request.Data.Description,
		Level:       request.Data.Level,
		Tags:        utils.SpliceStringWithX(request.Data.Tags, "#"),
		TestCase:    request.Data.TestCase,
	}
	result := global.DBInstance.Where("id = ?", question.ID)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "GetQuestionData error: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "method UpdateQuestionData error: %v", "not found")
	}
	result = global.DBInstance.Updates(question)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "method UpdateQuestionData error: %v", result.Error)
	}

	// 更新测试案例
	updateTestCases(request.Data.Id, request.Data.TestCase)

	return &empty.Empty{}, nil
}

func (receiver *DBServiceServer) DeleteQuestionData(ctx context.Context, request *pb.DeleteQuestionRequest) (*empty.Empty, error) {
	var question *models.Question
	result := global.DBInstance.Where("id = ?", request.Id).Find(&question)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "GetQuestionData error: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "GetQuestionData error: %v", "not found")
	}
	// 软删除
	result = global.DBInstance.Delete(question)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "method DeleteQuestionData error: %v", result.Error)
	}
	// 永久删除
	// result = global.DBInstance.Unscoped().Delete(&user)
	return &empty.Empty{}, nil
}

func (receiver *DBServiceServer) GetQuestionList(ctx context.Context, request *pb.GetQuestionListRequest) (*pb.GetQuestionListResponse, error) {
	var pageSize = 10
	rsp := &pb.GetQuestionListResponse{}

	var questionList []models.Question
	var count int64 = 0
	result := global.DBInstance.Model(questionList).Count(&count)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query count failed")
	}
	rsp.Total = int32(count)

	// select id,title from question
	// where id>=cursor
	// order by id
	// limit 10;
	result = global.DBInstance.Select("id, title").Where("id >= ?", request.Cursor).Order("id").Limit(pageSize).Find(&questionList)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query questionList failed")
	}
	for _, question := range questionList {
		rsp.Data = append(rsp.Data, &pb.Question{
			Id:    question.ID,
			Title: question.Title,
		})
	}
	rsp.Cursor = request.Cursor + int32(result.RowsAffected) + 1
	return rsp, nil
}

func (receiver *DBServiceServer) QueryQuestionWithName(ctx context.Context, request *pb.QueryQuestionWithNameRequest) (*pb.QueryQuestionWithNameResponse, error) {
	var questionList []models.Question
	// select * from question
	// where title like '%name%';
	names := "%" + request.Name + "%"
	result := global.DBInstance.Where("name LINK ?", names).Find(&questionList)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query question list failed")
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "method QueryQuestionWithName not found")
	}

	var data []*pb.Question
	for _, question := range questionList {
		data = append(data, &pb.Question{
			Id:          question.ID,
			Title:       question.Title,
			Description: question.Description,
			Level:       question.Level,
			Tags:        utils.SplitStringWithX(question.Tags, "#"),
		})
	}
	return &pb.QueryQuestionWithNameResponse{
		Data: data,
	}, nil
}

func updateTestCases(subKey int64, value string) {
	var key = "QuestionTestCases"
	conn := global.RedisPoolInstance.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", key, subKey, value)
	if err != nil {
		logrus.Errorf("UpdateTestCases error: %s", err.Error())
		return
	}
}
