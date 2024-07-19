package handler

import (
	"context"
	"db-service/internal/models"
	pb "db-service/internal/proto"
	"db-service/services/dao/mysql"
	"db-service/services/dao/redis"
	"db-service/utils"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (receiver *DBServiceServer) GetProblemData(ctx context.Context, request *pb.GetProblemRequest) (*pb.GetProblemResponse, error) {
	db := mysql.DB

	var problem models.Problem
	result := db.Where("id = ?", request.Id).Find(&problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, QueryField
	}
	if result.RowsAffected == 0 {
		return nil, NotFound
	}

	data := &pb.Problem{
		Id:          problem.ID,
		CreateAt:    timestamppb.New(problem.CreateAt),
		Title:       problem.Title,
		Description: problem.Description,
		Level:       problem.Level,
		Tags:        utils.SplitStringWithX(problem.Tags, "#"),
		TestCase:    problem.TestCase,
		CpuLimit:    problem.CpuLimit,
		ClockLimit:  problem.ClockLimit,
		MemoryLimit: problem.MemoryLimit,
		ProcLimit:   problem.ProcLimit,
		CreateBy:    problem.CreateBy,
	}
	return &pb.GetProblemResponse{
		Data: data,
	}, nil
}

func (receiver *DBServiceServer) CreateProblemData(ctx context.Context, request *pb.CreateProblemRequest) (*pb.CreateProblemResponse, error) {
	db := mysql.DB
	problem := &models.Problem{
		Title:       request.Data.Title,
		Level:       request.Data.Level,
		Tags:        utils.SpliceStringWithX(request.Data.Tags, "#"),
		Description: request.Data.Description,
		TestCase:    request.Data.TestCase,
		CpuLimit:    request.Data.CpuLimit,
		ClockLimit:  request.Data.ClockLimit,
		MemoryLimit: request.Data.MemoryLimit,
		ProcLimit:   request.Data.ProcLimit,
		CreateBy:    request.Data.CreateBy,
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

	// 插入成功 -- 将热点数据写到缓存(测试用例，题目配置)
	cacheProblemHotData(problem)

	return &pb.CreateProblemResponse{
		Id: problem.ID,
	}, nil
}

func (receiver *DBServiceServer) UpdateProblemData(ctx context.Context, request *pb.UpdateProblemRequest) (*empty.Empty, error) {
	db := mysql.DB
	problem := &models.Problem{
		Title:       request.Data.Title,
		Level:       request.Data.Level,
		Tags:        utils.SpliceStringWithX(request.Data.Tags, "#"),
		Description: request.Data.Description,
		TestCase:    request.Data.TestCase,
		CpuLimit:    request.Data.CpuLimit,
		ClockLimit:  request.Data.ClockLimit,
		MemoryLimit: request.Data.MemoryLimit,
		ProcLimit:   request.Data.ProcLimit,
		CreateBy:    request.Data.CreateBy,
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
	cacheProblemHotData(problem)

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

// 缓存题目热点数据
func cacheProblemHotData(problem *models.Problem) {
	data := &models.ProblemHotData{
		TestCase:    problem.TestCase,
		ClockLimit:  problem.ClockLimit,
		CpuLimit:    problem.CpuLimit,
		MemoryLimit: problem.MemoryLimit,
		ProcLimit:   problem.ProcLimit,
		TimeLimit:   problem.TimeLimit,
	}
	bys, err := json.Marshal(data)
	if err != nil {
		logrus.Errorln(err.Error())
		return
	}
	err = redis.SetKVByHash(fmt.Sprintf("problem:%d", problem.ID), "hotData", string(bys))
	if err != nil {
		logrus.Errorln(err.Error())
		return
	}
}
