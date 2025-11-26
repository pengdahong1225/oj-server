package service

import (
	"context"

	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"oj-server/global"
	"oj-server/pkg/gPool"
	"oj-server/pkg/mq"
	"oj-server/pkg/proto/pb"
	"oj-server/svr/problem/internal/biz"
	"oj-server/svr/problem/internal/configs"
	"oj-server/svr/problem/internal/data"
	"oj-server/svr/problem/internal/model"
	"time"
)

// record服务
type RecordService struct {
	pb.UnimplementedRecordServiceServer
	uc *biz.RecordUseCase

	result_consumer *mq.Consumer // 判题结果消费者
}

func NewRecordService() *RecordService {
	repo, err := data.NewRecordRepo()
	if err != nil {
		logrus.Fatalf("NewProblemService failed, err:%s", err.Error())
	}

	mqCfg := configs.AppConf.MQCfg
	amqpClient := mq.NewClient(
		&mq.Options{
			Host:     mqCfg.Host,
			Port:     mqCfg.Port,
			User:     mqCfg.User,
			PassWord: mqCfg.PassWord,
			VHost:    mqCfg.VHost,
		},
	)

	return &RecordService{
		uc: biz.NewRecordUseCase(repo),
		result_consumer: &mq.Consumer{
			AmqpClient: amqpClient, // 注入client
			ExKind:     global.RabbitMqExchangeKind,
			ExName:     global.RabbitMqExchangeName,
			QueName:    global.RabbitMqJudgeResultQueue,
			RoutingKey: global.RabbitMqJudgeResultKey,
			CTag:       "", // 消费者标签，用于区别不同的消费者
		},
	}
}

func (ps *RecordService) ConsumeJudgeResult() {
	deliveries := ps.result_consumer.Consume()
	if deliveries == nil {
		logrus.Errorf("获取deliveries失败")
		return
	}
	defer ps.result_consumer.Close()

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
			_ = gPool.Instance().Submit(func() {
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
	key := fmt.Sprintf("%s:%d", global.UserLockPrefix, result.Uid)
	err := ps.uc.UnLock(key)
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
	if err = ps.uc.UpdateSubmitRecord(result.TaskId, record, result.Level); err != nil {
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
	lb_list, err := ps.uc.QueryMonthAccomplishLeaderboard(200, time.Now().Format("2006-01"))
	if err != nil {
		logrus.Errorf("查询排行榜数据失败, err:%s", err.Error())
		return err
	}
	logrus.Debugf("获取月榜数据:%+v", lb_list)
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
	logrus.Debugf("获取日榜数据:%+v", lb_list)
	// 写入redis
	return ps.uc.SynchronizeLeaderboard(lb_list, global.GetDailyLeaderboardKey(), global.DailyLeaderboardTTL)
}

// 分页查询用户的提交记录
func (ps *RecordService) GetSubmitRecordList(ctx context.Context, in *pb.GetSubmitRecordListRequest) (*pb.GetSubmitRecordListResponse, error) {
	count, records, err := ps.uc.QuerySubmitRecordList(in.Uid, int(in.Page), int(in.PageSize))
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

func (ps *RecordService) QueryJudgeResult(ctx context.Context, in *pb.QueryJudgeResultRequest) (*pb.QueryJudgeResultResponse, error) {
	result, err := ps.uc.QueryJudgeResult(in.TaskId)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, err
	}
	return &pb.QueryJudgeResultResponse{
		Accepted: result.Accepted,
		Message:  result.Message,
	}, nil
}
