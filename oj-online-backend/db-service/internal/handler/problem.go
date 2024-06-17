package handler

import (
	"context"
	"db-service/internal/models"
	pb "db-service/internal/proto"
	"db-service/services/dao/mysql"
	"db-service/services/dao/redis"
	"db-service/utils"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetProblemData
// 获取题目信息
func (receiver *DBServiceServer) GetProblemData(ctx context.Context, request *pb.GetProblemRequest) (*pb.GetProblemResponse, error) {
	db := mysql.DB

	/**
	select problem.*, user_info.nickname
	from problem
	Left JOIN user_info on problem.create_by=user_info.id
	where problem.id = 1;
	*/
	var problemDataResult models.ProblemDataResult
	// result := db.Where("id = ?", request.Id).Find(&Problem)
	result := db.Table("problem").
		Select("problem.*, user_info.nickname").
		Joins("Left JOIN user_info on problem.create_by=user_info.id").
		Where("problem.id = ?", request.Id).
		Scan(&problemDataResult)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, QueryField
	}
	if result.RowsAffected == 0 {
		return nil, NotFound
	}

	data := &pb.Problem{
		Id:          problemDataResult.ID,
		CreateAt:    timestamppb.New(problemDataResult.CreateAt),
		Title:       problemDataResult.Title,
		Description: problemDataResult.Description,
		Level:       problemDataResult.Level,
		Tags:        utils.SplitStringWithX(problemDataResult.Tags, "#"),
		TestCase:    problemDataResult.TestCase,
		TimeLimit:   problemDataResult.TimeLimit,
		MemoryLimit: problemDataResult.MemoryLimit,
		IoMode:      problemDataResult.IoMode,
		CreateBy:    problemDataResult.CreateUserNickName,
	}
	return &pb.GetProblemResponse{
		Data: data,
	}, nil
}

func (receiver *DBServiceServer) CreateProblemData(ctx context.Context, request *pb.CreateProblemRequest) (*pb.CreateProblemResponse, error) {
	db := mysql.DB
	problem := &models.Problem{
		Title:       request.Data.Title,
		Description: request.Data.Description,
		Level:       request.Data.Level,
		Tags:        utils.SpliceStringWithX(request.Data.Tags, "#"),
		TestCase:    request.Data.TestCase,
	}
	result := db.Where("title = ?", problem.Title)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, QueryField
	}
	if result.RowsAffected > 0 {
		return nil, AlreadyExists
	}

	result = db.Create(problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, InsertFailed
	}

	// 插入成功 -- 将测试案例更新到redis
	updateTestCases(request.Data.Id, request.Data.TestCase)

	return &pb.CreateProblemResponse{
		Id: problem.ID,
	}, nil
}

func (receiver *DBServiceServer) UpdateProblemData(ctx context.Context, request *pb.UpdateProblemRequest) (*empty.Empty, error) {
	db := mysql.DB
	problem := &models.Problem{
		Title:       request.Data.Title,
		Description: request.Data.Description,
		Level:       request.Data.Level,
		Tags:        utils.SpliceStringWithX(request.Data.Tags, "#"),
		TestCase:    request.Data.TestCase,
	}
	result := db.Where("title = ?", problem.Title)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, QueryField
	}
	if result.RowsAffected == 0 {
		return nil, NotFound
	}

	// 更新
	result = db.Updates(problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, UpdateFailed
	}

	// 更新测试案例
	updateTestCases(request.Data.Id, request.Data.TestCase)

	return &empty.Empty{}, nil
}

func (receiver *DBServiceServer) DeleteProblemData(ctx context.Context, request *pb.DeleteProblemRequest) (*empty.Empty, error) {
	db := mysql.DB
	var problem *models.Problem
	result := db.Where("id = ?", request.Id).Find(&problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, QueryField
	}
	if result.RowsAffected == 0 {
		return nil, NotFound
	}

	// 软删除
	result = db.Delete(problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, DeleteFailed
	}
	// 永久删除
	// result = db.Unscoped().Delete(&user)
	return &empty.Empty{}, nil
}

// GetProblemList 题库列表
// 游标分页，查询id，title，level，tags
func (receiver *DBServiceServer) GetProblemList(ctx context.Context, request *pb.GetProblemListRequest) (*pb.GetProblemListResponse, error) {
	db := mysql.DB
	var pageSize = 10
	rsp := &pb.GetProblemListResponse{}

	var problemList []models.Problem
	var count int64 = 0
	result := db.Model(problemList).Count(&count)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, QueryField
	}
	rsp.Total = int32(count)

	// select id,title,level,tags from Problem
	// where id>=cursor
	// order by id
	// limit 10;
	result = db.Select("id,title,level,tags").Where("id >= ?", request.Cursor).Order("id").Limit(pageSize).Find(&problemList)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, QueryField
	}
	for _, Problem := range problemList {
		rsp.Data = append(rsp.Data, &pb.Problem{
			Id:    Problem.ID,
			Title: Problem.Title,
			Level: Problem.Level,
			Tags:  utils.SplitStringWithX(Problem.Tags, "#"),
		})
	}
	rsp.Cursor = request.Cursor + int32(result.RowsAffected) + 1
	return rsp, nil
}

// QueryProblemWithName 根据题目名查询题目
// 模糊查询
func (receiver *DBServiceServer) QueryProblemWithName(ctx context.Context, request *pb.QueryProblemWithNameRequest) (*pb.QueryProblemWithNameResponse, error) {
	db := mysql.DB
	var problemList []models.Problem
	// select * from Problem
	// where title like '%name%';
	names := "%" + request.Name + "%"
	result := db.Where("name LINK ?", names).Find(&problemList)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, QueryField
	}
	if result.RowsAffected == 0 {
		return nil, NotFound
	}

	var data []*pb.Problem
	for _, Problem := range problemList {
		data = append(data, &pb.Problem{
			Id:          Problem.ID,
			Title:       Problem.Title,
			Level:       Problem.Level,
			Tags:        utils.SplitStringWithX(Problem.Tags, "#"),
			Description: Problem.Description,
			TestCase:    Problem.TestCase,
		})
	}
	return &pb.QueryProblemWithNameResponse{
		Data: data,
	}, nil
}

// 测试用例要插入到redis
func updateTestCases(field int64, value string) {
	var key = "ProblemTestCases"
	conn := redis.NewConn()
	defer conn.Close()

	_, err := conn.Do("HSET", key, field, value)
	if err != nil {
		logrus.Errorf("UpdateTestCases error: %s", err.Error())
		return
	}
}
