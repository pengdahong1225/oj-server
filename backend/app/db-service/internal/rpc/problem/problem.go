package problem

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/rpc"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/svc/mysql"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/svc/redis"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/utils"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	problem := &mysql.Problem{
		Title:       request.Data.Title,
		Level:       request.Data.Level,
		Tags:        utils.SpliceStringWithX(request.Data.Tags, "#"),
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
	//if result.RowsAffected > 0 {
	//	return nil, AlreadyExists
	//}

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

	var problem mysql.Problem
	result := db.Where("id = ?", request.Id).Find(&problem)
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

	data := &pb.Problem{
		Id:          problem.ID,
		CreateAt:    timestamppb.New(problem.CreateAt),
		Title:       problem.Title,
		Description: problem.Description,
		Level:       problem.Level,
		Tags:        utils.SplitStringWithX(problem.Tags, "#"),
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

// GetProblemList 题库列表
// 游标分页，查询id，title，level，tags
func (receiver *ProblemServer) GetProblemList(ctx context.Context, request *pb.GetProblemListRequest) (*pb.GetProblemListResponse, error) {
	db := mysql.Instance()
	var pageSize = 10
	rsp := &pb.GetProblemListResponse{}

	var problemList []mysql.Problem
	var count int64 = 0
	result := db.Model(problemList).Count(&count)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc.QueryFailed
	}
	rsp.Total = int32(count)

	// select id,title,level,tags from Problem
	// where id>=cursor
	// order by id
	// limit 10;
	result = db.Select("id,title,level,tags").Where("id >= ?", request.Cursor).Order("id").Limit(pageSize).Find(&problemList)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc.QueryFailed
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
func (receiver *ProblemServer) QueryProblemWithName(ctx context.Context, request *pb.QueryProblemWithNameRequest) (*pb.QueryProblemWithNameResponse, error) {
	db := mysql.Instance()
	var problemList []mysql.Problem
	// select * from problem
	// where title like '%name%';
	names := "%" + request.Name + "%"
	result := db.Where("name LINK ?", names).Find(&problemList)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return nil, rpc.QueryFailed
	}
	if result.RowsAffected == 0 {
		return nil, rpc.NotFound
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
