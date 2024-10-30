package problem

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/rpc"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/svc/mysql"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/svc/redis"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"strconv"
)

type ProblemServer struct {
	pb.UnimplementedProblemServiceServer
}

func (receiver *ProblemServer) UpdateProblemData(ctx context.Context, request *pb.UpdateProblemRequest) (*pb.UpdateProblemResponse, error) {
	db := mysql.Instance()

	config, err := proto.Marshal(request.Data.Config)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, rpc.InsertFailed
	}
	tags, err := json.Marshal(request.Data.Tags)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, rpc.InsertFailed
	}

	problem := &mysql.Problem{
		Title:       request.Data.Title,
		Level:       request.Data.Level,
		Tags:        tags,
		Description: request.Data.Description,
		CreateBy:    request.Data.CreateBy,
		Config:      config,
	}
	result := db.Where("title = ?", problem.Title)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc.QueryFailed
	}

	// 不存在就新增，存在就修改
	if result.RowsAffected == 0 {
		result = db.Create(problem)
		if result.Error != nil {
			logrus.Errorln(result.Error.Error())
			return nil, rpc.InsertFailed
		}
	} else {
		result = db.Updates(problem)
		if result.Error != nil {
			logrus.Errorln(result.Error.Error())
			return nil, rpc.UpdateFailed
		}
	}

	// 缓存题目配置
	if err = redis.CacheProblemConfig(problem.ID, config); err != nil {
		logrus.Errorln(err.Error())
	}

	return &pb.UpdateProblemResponse{
		Id: problem.ID,
	}, nil
}

func (receiver *ProblemServer) GetProblemData(ctx context.Context, request *pb.GetProblemRequest) (*pb.GetProblemResponse, error) {
	db := mysql.Instance()

	/*
		SELECT * FROM `problem`
		where id = ?;
	*/
	var problem mysql.Problem

	fmt.Println(db.Debug().First(&problem, "id = ?", request.Id).Statement.SQL.String())

	result := db.Where("id = ?", request.Id).First(&problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc.QueryFailed
	}
	if result.RowsAffected == 0 {
		return nil, rpc.NotFound
	}

	config := &pb.ProblemConfig{}
	err := proto.Unmarshal(problem.Config, config)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, rpc.QueryFailed
	}

	var tags []string
	err = json.Unmarshal(problem.Tags, &tags)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, rpc.QueryFailed
	}

	data := &pb.Problem{
		Id:          problem.ID,
		CreateAt:    strconv.FormatInt(problem.CreateAt.Unix(), 10),
		Title:       problem.Title,
		Description: problem.Description,
		Level:       problem.Level,
		Tags:        tags,
		CreateBy:    problem.CreateBy,
		Config:      config,
	}

	return &pb.GetProblemResponse{
		Data: data,
	}, nil
}

func (receiver *ProblemServer) DeleteProblemData(ctx context.Context, request *pb.DeleteProblemRequest) (*empty.Empty, error) {
	db := mysql.Instance()
	var problem *mysql.Problem
	result := db.Where("id = ?", request.Id).Find(&problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc.QueryFailed
	}
	if result.RowsAffected == 0 {
		return nil, rpc.NotFound
	}

	// 软删除
	result = db.Delete(problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc.DeleteFailed
	}
	// 永久删除
	// result = db.Unscoped().Delete(&user)
	return &empty.Empty{}, nil
}

// GetProblemList 分页查询题库列表
// 查询{id，title，level，tags}
// @page 页码
// @page_size 单页数量
// @keyword 关键字
// @tag 标签
func (receiver *ProblemServer) GetProblemList(ctx context.Context, request *pb.GetProblemListRequest) (*pb.GetProblemListResponse, error) {
	db := mysql.Instance()
	rsp := &pb.GetProblemListResponse{}
	name := "%" + request.Keyword + "%"
	offSet := int((request.Page - 1) * request.PageSize)
	query := fmt.Sprintf(`JSON_CONTAINS(tags, '"%s"')`, request.Tag)
	logrus.Debugln("query conditions: %s", query)

	/*
		select COUNT(*) AS count from problem
		where title like '%name%' AND JSON_CONTAINS(tags, '"哈希表"');
	*/
	var count int64 = 0
	result := db.Model(&mysql.Problem{}).Where("title LIKE ?", name).Where(query).Count(&count)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc.QueryFailed
	}
	rsp.Total = int32(count)

	/*
		select id,title,level,tags from problem
		where title like '%name%' AND JSON_CONTAINS(tags, '"哈希表')
		order by id
		offset off_set
		limit page_size;
	*/
	var problemList []mysql.Problem
	result = db.Select("id,title,level,tags").Where("title LIKE ?", name).Where(query).Order("id").Offset(offSet).Limit(int(request.PageSize)).Find(&problemList)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc.QueryFailed
	}
	for _, problem := range problemList {
		p := &pb.Problem{
			Id:    problem.ID,
			Title: problem.Title,
			Level: problem.Level,
		}
		json.Unmarshal(problem.Tags, &p.Tags)
		rsp.Data = append(rsp.Data, p)
	}
	return rsp, nil
}

// GetProblemHotData 读取题目热点数据
// 获取后更新到缓存
func (receiver *ProblemServer) GetProblemHotData(ctx context.Context, request *pb.GetProblemHotDataRequest) (*pb.GetProblemHotDataResponse, error) {
	db := mysql.Instance()
	var problem mysql.Problem
	// select config
	// from problem
	// where id = ?
	result := db.Select("config").Where("id = ?", request.ProblemId).Find(&problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc.QueryFailed
	}
	if result.RowsAffected == 0 {
		return nil, rpc.NotFound
	}

	err := redis.CacheProblemConfig(problem.ID, problem.Config)
	if err != nil {
		logrus.Errorln(err.Error())
	}

	response := &pb.GetProblemHotDataResponse{Data: string(problem.Config)}
	return response, nil
}
