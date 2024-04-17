package handler

import (
	"context"
	"db-service/internal/models"
	pb2 "db-service/internal/proto"
	"db-service/services/dao/mysql"
	"db-service/services/dao/redis"
	"db-service/utils"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (receiver *DBServiceServer) GetQuestionData(ctx context.Context, request *pb2.GetQuestionRequest) (*pb2.GetQuestionResponse, error) {
	db := mysql.DB
	var question models.Question
	result := db.Where("id = ?", request.Id).Find(&question)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "GetQuestionData error: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "GetQuestionData[%d] error: %v", request.Id, "not found")
	}

	data := &pb2.Question{
		Id:          question.ID,
		CreateAt:    timestamppb.New(question.CreateAt),
		Title:       question.Title,
		Level:       question.Level,
		Tags:        utils.SplitStringWithX(question.Tags, "#"),
		Description: question.Description,
		TestCase:    question.TestCase,
		Template:    question.Template,
	}
	return &pb2.GetQuestionResponse{
		Data: data,
	}, nil
}

func (receiver *DBServiceServer) CreateQuestionData(ctx context.Context, request *pb2.CreateQuestionRequest) (*pb2.CreateQuestionResponse, error) {
	db := mysql.DB
	question := &models.Question{
		Title:       request.Data.Title,
		Level:       request.Data.Level,
		Tags:        utils.SpliceStringWithX(request.Data.Tags, "#"),
		Description: request.Data.Description,
		TestCase:    request.Data.TestCase,
		Template:    request.Data.Template,
	}
	result := db.Where("id = ?", question.ID)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "GetQuestionData error: %v", result.Error)
	}
	if result.RowsAffected > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "method CreateQuestionData error: %v", "already exists")
	}

	result = db.Create(question)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "method CreateQuestionData error: %v", result.Error)
	}

	// 插入成功 -- 将测试案例插入到redis
	updateTestCases(request.Data.Id, request.Data.TestCase)

	return &pb2.CreateQuestionResponse{
		Id: question.ID,
	}, nil
}

func (receiver *DBServiceServer) UpdateQuestionData(ctx context.Context, request *pb2.UpdateQuestionRequest) (*empty.Empty, error) {
	db := mysql.DB
	question := &models.Question{
		Title:       request.Data.Title,
		Level:       request.Data.Level,
		Tags:        utils.SpliceStringWithX(request.Data.Tags, "#"),
		Description: request.Data.Description,
		TestCase:    request.Data.TestCase,
		Template:    request.Data.Template,
	}
	result := db.Where("id = ?", question.ID)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "GetQuestionData error: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "method UpdateQuestionData error: %v", "not found")
	}
	result = db.Updates(question)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "method UpdateQuestionData error: %v", result.Error)
	}

	// 更新测试案例
	updateTestCases(request.Data.Id, request.Data.TestCase)

	return &empty.Empty{}, nil
}

func (receiver *DBServiceServer) DeleteQuestionData(ctx context.Context, request *pb2.DeleteQuestionRequest) (*empty.Empty, error) {
	db := mysql.DB
	var question *models.Question
	result := db.Where("id = ?", request.Id).Find(&question)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "GetQuestionData error: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "GetQuestionData error: %v", "not found")
	}
	// 软删除
	result = db.Delete(question)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "method DeleteQuestionData error: %v", result.Error)
	}
	// 永久删除
	// result = db.Unscoped().Delete(&user)
	return &empty.Empty{}, nil
}

// GetQuestionList 题库列表
// 游标分页，查询id，title，level，tags
func (receiver *DBServiceServer) GetQuestionList(ctx context.Context, request *pb2.GetQuestionListRequest) (*pb2.GetQuestionListResponse, error) {
	db := mysql.DB
	var pageSize = 10
	rsp := &pb2.GetQuestionListResponse{}

	var questionList []models.Question
	var count int64 = 0
	result := db.Model(questionList).Count(&count)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query count failed")
	}
	rsp.Total = int32(count)

	// select id,title,level,tags from question
	// where id>=cursor
	// order by id
	// limit 10;
	result = db.Select("id,title,level,tags").Where("id >= ?", request.Cursor).Order("id").Limit(pageSize).Find(&questionList)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query questionList failed")
	}
	for _, question := range questionList {
		rsp.Data = append(rsp.Data, &pb2.Question{
			Id:    question.ID,
			Title: question.Title,
			Level: question.Level,
			Tags:  utils.SplitStringWithX(question.Tags, "#"),
		})
	}
	rsp.Cursor = request.Cursor + int32(result.RowsAffected) + 1
	return rsp, nil
}

// QueryQuestionWithName 根据题目名查询题目
// 模糊查询
func (receiver *DBServiceServer) QueryQuestionWithName(ctx context.Context, request *pb2.QueryQuestionWithNameRequest) (*pb2.QueryQuestionWithNameResponse, error) {
	db := mysql.DB
	var questionList []models.Question
	// select * from question
	// where title like '%name%';
	names := "%" + request.Name + "%"
	result := db.Where("name LINK ?", names).Find(&questionList)
	if result.Error != nil {
		logrus.Debugln(result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query question list failed")
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "method QueryQuestionWithName not found")
	}

	var data []*pb2.Question
	for _, question := range questionList {
		data = append(data, &pb2.Question{
			Id:          question.ID,
			Title:       question.Title,
			Level:       question.Level,
			Tags:        utils.SplitStringWithX(question.Tags, "#"),
			Description: question.Description,
			TestCase:    question.TestCase,
			Template:    question.Template,
		})
	}
	return &pb2.QueryQuestionWithNameResponse{
		Data: data,
	}, nil
}

// 测试用例要插入到redis
func updateTestCases(subKey int64, value string) {
	var key = "QuestionTestCases"
	conn := redis.NewConn()
	defer conn.Close()

	_, err := conn.Do("HSET", key, subKey, value)
	if err != nil {
		logrus.Errorf("UpdateTestCases error: %s", err.Error())
		return
	}
}
