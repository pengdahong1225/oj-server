package handler

import (
	"context"
	pb "db-service/internal/proto"
	mysql2 "db-service/services/mysql"
	"db-service/services/redis"
	"db-service/utils"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (receiver *DBServiceServer) CreateProblemData(ctx context.Context, request *pb.CreateProblemRequest) (*pb.CreateProblemResponse, error) {
	db := mysql2.DB

	compileConfig := mysql2.ProblemConfig{
		ClockLimit:  request.Data.CompileConfig.ClockLimit,
		CpuLimit:    request.Data.CompileConfig.CpuLimit,
		MemoryLimit: request.Data.CompileConfig.MemoryLimit,
		ProcLimit:   request.Data.CompileConfig.ProcLimit,
	}
	runConfig := mysql2.ProblemConfig{
		ClockLimit:  request.Data.RunConfig.ClockLimit,
		CpuLimit:    request.Data.RunConfig.CpuLimit,
		MemoryLimit: request.Data.RunConfig.MemoryLimit,
		ProcLimit:   request.Data.RunConfig.ProcLimit,
	}
	var testCases []mysql2.TestCase
	for _, test := range request.Data.TestCases {
		testCases = append(testCases, mysql2.TestCase{
			Input:  test.Input,
			Output: test.Output,
		})
	}
	cbys, err := json.Marshal(&compileConfig)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, InsertFailed
	}
	rbys, err := json.Marshal(&runConfig)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, InsertFailed
	}
	tbys, err := json.Marshal(testCases)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, InsertFailed
	}

	problem := &mysql2.Problem{
		Title:         request.Data.Title,
		Level:         request.Data.Level,
		Tags:          utils.SpliceStringWithX(request.Data.Tags, "#"),
		Description:   request.Data.Description,
		CreateBy:      request.Data.CreateBy,
		TestCase:      string(tbys),
		CompileConfig: string(cbys),
		RunConfig:     string(rbys),
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

	// 将热点数据写到缓存(测试用例，题目配置)
	cacheProblemHotData(problem)

	return &pb.CreateProblemResponse{
		Id: problem.ID,
	}, nil
}

func (receiver *DBServiceServer) GetProblemData(ctx context.Context, request *pb.GetProblemRequest) (*pb.GetProblemResponse, error) {
	db := mysql2.DB

	var problem mysql2.Problem
	result := db.Where("id = ?", request.Id).Find(&problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, QueryField
	}
	if result.RowsAffected == 0 {
		return nil, NotFound
	}

	var compileConfig mysql2.ProblemConfig
	var runConfig mysql2.ProblemConfig
	var testCases []mysql2.TestCase
	err := json.Unmarshal([]byte(problem.CompileConfig), &compileConfig)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, QueryField
	}
	err = json.Unmarshal([]byte(problem.RunConfig), &runConfig)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, QueryField
	}
	err = json.Unmarshal([]byte(problem.TestCase), &testCases)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, QueryField
	}

	data := &pb.Problem{
		Id:          problem.ID,
		CreateAt:    timestamppb.New(problem.CreateAt),
		Title:       problem.Title,
		Description: problem.Description,
		Level:       problem.Level,
		Tags:        utils.SplitStringWithX(problem.Tags, "#"),
		CreateBy:    problem.CreateBy,
		CompileConfig: &pb.ProblemConfig{
			ClockLimit:  compileConfig.ClockLimit,
			CpuLimit:    compileConfig.CpuLimit,
			MemoryLimit: compileConfig.MemoryLimit,
			ProcLimit:   compileConfig.ProcLimit,
		},
		RunConfig: &pb.ProblemConfig{
			ClockLimit:  runConfig.ClockLimit,
			CpuLimit:    runConfig.CpuLimit,
			MemoryLimit: runConfig.MemoryLimit,
			ProcLimit:   runConfig.ProcLimit,
		},
	}
	for _, test := range testCases {
		data.TestCases = append(data.TestCases, &pb.TestCase{
			Input:  test.Input,
			Output: test.Output,
		})
	}

	return &pb.GetProblemResponse{
		Data: data,
	}, nil
}

func (receiver *DBServiceServer) UpdateProblemData(ctx context.Context, request *pb.UpdateProblemRequest) (*empty.Empty, error) {
	db := mysql2.DB
	compileConfig := mysql2.ProblemConfig{
		ClockLimit:  request.Data.CompileConfig.ClockLimit,
		CpuLimit:    request.Data.CompileConfig.CpuLimit,
		MemoryLimit: request.Data.CompileConfig.MemoryLimit,
		ProcLimit:   request.Data.CompileConfig.ProcLimit,
	}
	runConfig := mysql2.ProblemConfig{
		ClockLimit:  request.Data.RunConfig.ClockLimit,
		CpuLimit:    request.Data.RunConfig.CpuLimit,
		MemoryLimit: request.Data.RunConfig.MemoryLimit,
		ProcLimit:   request.Data.RunConfig.ProcLimit,
	}
	var testCases []mysql2.TestCase
	for _, test := range request.Data.TestCases {
		testCases = append(testCases, mysql2.TestCase{
			Input:  test.Input,
			Output: test.Output,
		})
	}
	cbys, err := json.Marshal(&compileConfig)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, InsertFailed
	}
	rbys, err := json.Marshal(&runConfig)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, InsertFailed
	}
	tbys, err := json.Marshal(testCases)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, InsertFailed
	}

	problem := &mysql2.Problem{
		Title:         request.Data.Title,
		Level:         request.Data.Level,
		Tags:          utils.SpliceStringWithX(request.Data.Tags, "#"),
		Description:   request.Data.Description,
		CreateBy:      request.Data.CreateBy,
		TestCase:      string(tbys),
		CompileConfig: string(cbys),
		RunConfig:     string(rbys),
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
	db := mysql2.DB
	var problem *mysql2.Problem
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
	db := mysql2.DB
	var pageSize = 10
	rsp := &pb.GetProblemListResponse{}

	var problemList []mysql2.Problem
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
	db := mysql2.DB
	var problemList []mysql2.Problem
	// select * from problem
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
		})
	}
	return &pb.QueryProblemWithNameResponse{
		Data: data,
	}, nil
}

// GetProblemHotData 读取题目热点数据
// 获取后更新到缓存
func (receiver *DBServiceServer) GetProblemHotData(ctx context.Context, request *pb.GetProblemHotDataRequest) (*pb.GetProblemHotDataResponse, error) {
	db := mysql2.DB
	var problem mysql2.Problem
	// select test_case, compile_config, run_config
	// from problem
	// where id = ?
	result := db.Where("id = ?", request.ProblemId).Find(&problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, QueryField
	}
	if result.RowsAffected == 0 {
		return nil, NotFound
	}

	str := cacheProblemHotData(&problem)

	response := &pb.GetProblemHotDataResponse{Data: str}
	return response, nil
}

// 缓存题目热点数据
func cacheProblemHotData(problem *mysql2.Problem) string {
	data := &mysql2.ProblemHotData{
		TestCase:      problem.TestCase,
		CompileConfig: problem.CompileConfig,
		RunConfig:     problem.RunConfig,
	}
	bys, err := json.Marshal(data)
	if err != nil {
		logrus.Errorln(err.Error())
		return ""
	}
	err = redis.SetKVByHash(fmt.Sprintf("problem:%d", problem.ID), "hotData", string(bys))
	if err != nil {
		logrus.Errorln(err.Error())
		return ""
	}
	return string(bys)
}
