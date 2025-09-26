package service

import (
	"context"

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

// 系统启动时，先全量同步一次
func (ps *RecordService) SyncLeaderboardByScheduled() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("update leaderboard panic: %v", err)
		}
	}()
	if err := ps.syncMonthLeaderboard(); err != nil {
		logrus.Errorf("全量同步月榜失败, err:%s", err.Error())
	}
	if err := ps.syncMonthLeaderboard(); err != nil {
		logrus.Errorf("全量同步日榜失败, err:%s", err.Error())
	}

	logrus.Infof("排行榜建立成功")

	// todo 定时补偿 -- 防止漏更新

}

// 默认只维护200条数据
func (ps *RecordService) syncMonthLeaderboard() error {
	// 从数据库中获取数据
	lb_list, err := ps.uc.QueryMonthAccomplishLeaderboard(200, time.Now().Format("2006-01"))
	if err != nil {
		logrus.Errorf("查询排行榜数据失败, err:%s", err.Error())
		return err
	}
	// 写入redis
	return ps.uc.SynchronizeLeaderboard(lb_list, global.GetMonthLeaderboardKey(), global.MonthLeaderboardTTL)
}
func (ps *RecordService) syncDailyLeaderboard() error {
	// 从数据库中获取数据
	lb_list, err := ps.uc.QueryDailyAccomplishLeaderboard(200)
	if err != nil {
		logrus.Errorf("查询排行榜数据失败, err:%s", err.Error())
		return err
	}
	// 写入redis
	return ps.uc.SynchronizeLeaderboard(lb_list, global.GetDailyLeaderboardKey(), global.DailyLeaderboardTTL)
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
