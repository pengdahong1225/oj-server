package service

import (
	"context"

	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"oj-server/global"
	"oj-server/proto/pb"
	"oj-server/svr/problem/internal/biz"
	"oj-server/svr/problem/internal/data"
	"time"
)

// record服务
type RecordService struct {
	pb.UnimplementedRecordServiceServer
	uc *biz.RecordUseCase
}

func NewRecordService() *RecordService {
	repo, err := data.NewRecordRepo()
	if err != nil {
		logrus.Fatalf("NewProblemService failed, err:%s", err.Error())
	}

	return &RecordService{
		uc: biz.NewRecordUseCase(repo),
	}
}

// 排行榜定时维护
func (ps *RecordService) UpdateLeaderboardByScheduled() {
	// 系统启动时，先全量同步一次
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("update leaderboard panic: %v", err)
			}
		}()
		if err := ps.syncLeaderboard(); err != nil {
			logrus.Errorf("同步排行榜失败, err:%s", err.Error())
		}
	}()
	randListUpdateTicker := time.NewTicker(global.LeaderboardTTL)
	defer randListUpdateTicker.Stop()
	for range randListUpdateTicker.C {
		func() {
			defer func() {
				if err := recover(); err != nil {
					logrus.Errorf("update leaderboard panic: %v", err)
				}
			}()
			logrus.Infof("<----------------定时同步排行榜---------------->")
			// 检查是否需要更新
			lastUpdated, err := ps.uc.QueryLeaderboardLastUpdate()
			if err != nil {
				switch {
				case errors.Is(err, redis.Nil):
					lastUpdated = time.Now().Unix()
					if err = ps.uc.UpdateLeaderboardLastUpdate(lastUpdated); err != nil {
						logrus.Errorf("更新排行榜最后更新时间失败, err:%s", err.Error())
						return
					}
				default:
					logrus.Errorf("查询排行榜最后更新时间失败, err:%s", err.Error())
					return
				}
			}
			// 如果未超过更新间隔, 跳过
			if time.Now().Unix()-lastUpdated < int64(global.LeaderboardTTL.Seconds()) {
				return
			}
			// 同步排行榜
			if err = ps.syncLeaderboard(); err != nil {
				logrus.Errorf("同步排行榜失败, err:%s", err.Error())
			}
			// 更新排行榜最后更新时间
			if err = ps.uc.UpdateLeaderboardLastUpdate(time.Now().Unix()); err != nil {
				logrus.Errorf("更新排行榜最后更新时间失败, err:%s", err.Error())
				return
			}
		}()
	}
}
func (ps *RecordService) syncLeaderboard() error {
	// todo 从数据库中获取数据

	// todo 使用pipe批量操作redis
}

// 分页查询用户的提交记录
func (ps *RecordService) GetSubmitRecordList(ctx context.Context, in *pb.GetSubmitRecordListRequest) (*pb.GetSubmitRecordListResponse, error) {
	offSet := int((in.Page - 1) * in.PageSize)
	count, records, err := ps.uc.QuerySubmitRecordList(in.Uid, int(in.PageSize), offSet)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, err
	}

	resp := &pb.GetSubmitRecordListResponse{
		Total: int32(count),
		Data:  make([]*pb.SubmitRecord, len(records)),
	}
	for i, record := range records {
		resp.Data[i] = &pb.SubmitRecord{
			Id:          int64(record.ID),
			CreatedAt:   record.CreatedAt.Unix(),
			ProblemName: record.ProblemName,
			Lang:        record.Lang,
			Status:      record.Status,
		}
	}

	return resp, nil
}

// 获取提交记录数据
func (ps *RecordService) GetSubmitRecordData(ctx context.Context, in *pb.GetSubmitRecordRequest) (*pb.GetSubmitRecordResponse, error) {
	record, err := ps.uc.QuerySubmitRecord(in.Id)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, err
	}
	resp := &pb.GetSubmitRecordResponse{
		Data: record.Transform(),
	}

	return resp, nil
}
