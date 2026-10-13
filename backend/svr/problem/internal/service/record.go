package service

import (
	"context"

	"fmt"
	"oj-server/common/consts"
	"oj-server/common/proto/pb"
	"oj-server/pkg/go_pool"
	"oj-server/pkg/mq"
	"oj-server/svr/problem/internal/client"
	"oj-server/svr/problem/internal/config"
	"oj-server/svr/problem/internal/data"
	"oj-server/svr/problem/internal/model"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

// record服务
type RecordService struct {
	pb.UnimplementedRecordServiceServer
	repo           *data.RecordRepo
	resultConsumer *mq.Consumer
	clients        *client.Clients
}

func NewRecordService(clients *client.Clients) (*RecordService, error) {
	repo, err := data.NewRecordRepo()
	if err != nil {
		logrus.Errorf("NewProblemService failed, err:%s", err.Error())
		return nil, err
	}

	cfg := config.MyConfig.MQ
	amqpClient := mq.NewClient(
		&mq.Options{
			Host:     cfg.Host,
			Port:     cfg.Port,
			User:     cfg.User,
			PassWord: cfg.PassWord,
			VHost:    cfg.VHost,
		},
	)

	return &RecordService{
		repo: repo,
		resultConsumer: &mq.Consumer{
			AmqpClient: amqpClient, // 注入client
			ExKind:     consts.RabbitMqExchangeKind,
			ExName:     consts.RabbitMqExchangeName,
			QueName:    consts.RabbitMqJudgeResultQueue,
			RoutingKey: consts.RabbitMqJudgeResultKey,
			CTag:       "", // 消费者标签，用于区别不同的消费者
		},
		clients: clients,
	}, nil
}

func (ps *RecordService) ConsumeJudgeResult() {
	deliveries := ps.resultConsumer.Consume()
	if deliveries == nil {
		logrus.Errorf("获取deliveries失败")
		return
	}
	defer ps.resultConsumer.Close()

	for d := range deliveries {
		// 处理任务
		result := func(data []byte) bool {
			result := new(pb.JudgeResult)
			err := proto.Unmarshal(data, result)
			if err != nil {
				logrus.Errorln("解析judge task err：", err.Error())
				return false
			}
			// 异步处理
			_ = go_pool.Submit(func() {
				ps.handleJudgeResult(result)
			})
			return true
		}(d.Body)

		// 确认
		if result {
			_ = d.Ack(false)
		} else {
			_ = d.Reject(false)
		}
	}
}

func (ps *RecordService) handleJudgeResult(result *pb.JudgeResult) {
	// 释放锁
	key := fmt.Sprintf("%s:%d", consts.UserLockPrefix, result.Uid)
	err := ps.repo.UnLock(key)
	if err != nil {
		logrus.Errorf("释放锁失败, err:%s", err.Error())
	}

	// 更新数据库
	record := &model.SubmitRecord{
		UID:         result.Uid,
		UserName:    result.UserName,
		ProblemID:   result.ProblemId,
		ProblemName: result.ProblemName,
		Accepted:    result.Accepted,
		Message:     result.Message,
		Lang:        result.Lang,
		Code:        result.Code,
	}
	judgeResultStore := &pb.JudgeResultStore{
		Items: result.Items,
	}
	record.Result, err = proto.Marshal(judgeResultStore)
	if err != nil {
		logrus.Errorf("proto marshal err：%s", err.Error())
		return
	}
	if err = ps.repo.UpdateSubmitRecord(result.TaskId, record, result.Level); err != nil {
		logrus.Errorf("更新record失败, err:%s", err.Error())
	}
}

// 系统启动时，先全量同步一次
func (ps *RecordService) SyncLeaderboardByScheduled() {
	//defer func() {
	//	if err := recover(); err != nil {
	//		logrus.Errorf("update leaderboard panic: %v", err)
	//	}
	//}()
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
	lb_list, err := ps.repo.QueryMonthAccomplishLeaderboard(200, time.Now().Format("2006-01"))
	if err != nil {
		logrus.Errorf("查询排行榜数据失败, err:%s", err.Error())
		return err
	}
	logrus.Debugf("获取月榜数据:%+v", lb_list)
	// 写入redis
	return ps.repo.SynchronizeLeaderboard(lb_list, consts.GetMonthLeaderboardKey(), consts.MonthLeaderboardTTL)
}
func (ps *RecordService) syncDailyLeaderboard() error {
	// 从数据库中获取数据
	lb_list, err := ps.repo.QueryDailyAccomplishLeaderboard(200)
	if err != nil {
		logrus.Errorf("查询排行榜数据失败, err:%s", err.Error())
		return err
	}
	logrus.Debugf("获取日榜数据:%+v", lb_list)
	// 写入redis
	return ps.repo.SynchronizeLeaderboard(lb_list, consts.GetDailyLeaderboardKey(), consts.DailyLeaderboardTTL)
}

// 分页查询用户的提交记录
func (ps *RecordService) GetSubmitRecordList(ctx context.Context, in *pb.GetSubmitRecordListRequest) (*pb.GetSubmitRecordListResponse, error) {
	count, records, err := ps.repo.QuerySubmitRecordList(in.Uid, int(in.Page), int(in.PageSize))
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
			Id:          record.ID,
			CreatedAt:   record.CreatedAt.Unix(),
			Uid:         record.UID,
			UserName:    record.UserName,
			ProblemId:   record.ProblemID,
			ProblemName: record.ProblemName,
			Lang:        record.Lang,
			Accepted:    record.Accepted,
			Message:     record.Message,
		}
	}

	return resp, nil
}

// 获取提交记录数据
func (ps *RecordService) GetSubmitRecordData(ctx context.Context, in *pb.GetSubmitRecordRequest) (*pb.GetSubmitRecordResponse, error) {
	record, err := ps.repo.QuerySubmitRecord(in.Id)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, err
	}
	resp := &pb.GetSubmitRecordResponse{
		Data: record.Transform(),
	}

	return resp, nil
}

func (ps *RecordService) QueryJudgeResult(ctx context.Context, in *pb.QueryJudgeResultRequest) (*pb.QueryJudgeResultResponse, error) {
	result, err := ps.repo.QueryJudgeResult(in.TaskId)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, err
	}
	return &pb.QueryJudgeResultResponse{
		Accepted: result.Accepted,
		Message:  result.Message,
	}, nil
}
