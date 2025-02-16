package problem

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pengdahong1225/oj-server/backend/app/common/errs"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/cache"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/mysql"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"strconv"
)

type ProblemServer struct {
	pb.UnimplementedProblemServiceServer
}

func (receiver *ProblemServer) UpdateProblemData(ctx context.Context, request *pb.UpdateProblemRequest) (*pb.UpdateProblemResponse, error) {
	db := mysql.DBSession

	config, err := proto.Marshal(request.Data.Config)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errs.InsertFailed
	}
	tags, err := json.Marshal(request.Data.Tags)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errs.InsertFailed
	}

	problem := mysql.Problem{}
	result := db.Where("title = ?", request.Data.Title).Find(&problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, errs.QueryFailed
	}

	problem.Title = request.Data.Title
	problem.Level = request.Data.Level
	problem.Tags = tags
	problem.Description = request.Data.Description
	problem.CreateBy = request.Data.CreateBy
	problem.Config = config

	if result.RowsAffected == 0 {
		result = db.Create(&problem)
	} else {
		result = db.Updates(&problem)
	}
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, errs.UpdateFailed
	}

	// 缓存题目配置
	if err = cache.CacheProblemConfig(problem.ID, config); err != nil {
		logrus.Errorln(err.Error())
	}
	// 缓存标签
	if err = cache.UpdateTagList(request.Data.Tags); err != nil {
		logrus.Errorln(err.Error())
	}

	return &pb.UpdateProblemResponse{
		Id: problem.ID,
	}, nil
}

func (receiver *ProblemServer) GetProblemData(ctx context.Context, request *pb.GetProblemRequest) (*pb.GetProblemResponse, error) {
	db := mysql.DBSession

	/*
		SELECT * FROM `problem`
		where id = ?;
	*/
	var problem mysql.Problem

	result := db.Where("id = ?", request.Id).First(&problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, errs.QueryFailed
	}
	if result.RowsAffected == 0 {
		return nil, errs.NotFound
	}

	config := &pb.ProblemConfig{}
	err := proto.Unmarshal(problem.Config, config)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errs.QueryFailed
	}

	var tags []string
	err = json.Unmarshal(problem.Tags, &tags)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, errs.QueryFailed
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
	db := mysql.DBSession
	var problem *mysql.Problem
	result := db.Where("id = ?", request.Id).Find(&problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, errs.QueryFailed
	}
	if result.RowsAffected == 0 {
		return nil, errs.NotFound
	}

	// 软删除
	result = db.Delete(problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, errs.DeleteFailed
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
	db := mysql.DBSession
	rsp := &pb.GetProblemListResponse{}
	name := "%" + request.Keyword + "%"
	offSet := int((request.Page - 1) * request.PageSize)
	query := fmt.Sprintf(`JSON_CONTAINS(tags, '"%s"')`, request.Tag)
	logrus.Debugf("query conditions: %s\n", query)

	var result *gorm.DB
	/*
		select COUNT(*) AS count from problem
		where title like '%name%' AND JSON_CONTAINS(tags, '"哈希表"');
	*/
	var count int64 = 0
	if request.Tag == "" {
		result = db.Model(&mysql.Problem{}).Where("title LIKE ?", name).Count(&count)
	} else {
		result = db.Model(&mysql.Problem{}).Where("title LIKE ?", name).Where(query).Count(&count)
	}
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, errs.QueryFailed
	}
	rsp.Total = int32(count)
	logrus.Debugln("count:", count)

	/*
		select id,title,level,tags,create_at,create_by from problem
		where title like '%name%' AND JSON_CONTAINS(tags, '"哈希表')
		order by id
		offset off_set
		limit page_size;
	*/
	var problemList []mysql.Problem
	if request.Tag == "" {
		result = db.Select("id,title,level,tags,create_at,create_by").Where("title LIKE ?", name).Order("id").Offset(offSet).Limit(int(request.PageSize)).Find(&problemList)
	} else {
		result = db.Select("id,title,level,tags,create_at,create_by").Where("title LIKE ?", name).Where(query).Order("id").Offset(offSet).Limit(int(request.PageSize)).Find(&problemList)
	}

	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, errs.QueryFailed
	}
	for _, problem := range problemList {
		p := &pb.Problem{
			Id:       problem.ID,
			Title:    problem.Title,
			Level:    problem.Level,
			CreateAt: strconv.FormatInt(problem.CreateAt.Unix(), 10),
			CreateBy: problem.CreateBy,
		}
		err := json.Unmarshal(problem.Tags, &p.Tags)
		if err != nil {
			logrus.Errorln(err.Error())
			return nil, errs.QueryFailed
		}
		rsp.Data = append(rsp.Data, p)
	}
	return rsp, nil
}

// GetProblemHotData 读取题目热点数据
// 获取后更新到缓存
func (receiver *ProblemServer) GetProblemHotData(ctx context.Context, request *pb.GetProblemHotDataRequest) (*pb.GetProblemHotDataResponse, error) {
	db := mysql.DBSession
	var problem mysql.Problem
	// select config
	// from problem
	// where id = ?
	result := db.Select("config").Where("id = ?", request.ProblemId).Find(&problem)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, errs.QueryFailed
	}
	if result.RowsAffected == 0 {
		return nil, errs.NotFound
	}

	err := cache.CacheProblemConfig(problem.ID, problem.Config)
	if err != nil {
		logrus.Errorln(err.Error())
	}

	response := &pb.GetProblemHotDataResponse{Data: problem.Config}
	return response, nil
}
