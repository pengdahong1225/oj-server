package judge

import (
	"context"
	"encoding/json"
	"github.com/pengdahong1225/oj-server/backend/app/judge-service/internal/svc/cache"
	"github.com/pengdahong1225/oj-server/backend/app/judge-service/internal/types"
	"github.com/pengdahong1225/oj-server/backend/module/goroutinePool"
	"github.com/pengdahong1225/oj-server/backend/module/registry"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"sync"
	"time"
)

// Handle 判题服务入口
func Handle(form *pb.SubmitForm) {
	// 退出之后，需要将本次提交的状态置为UPStateExited
	defer func() {
		if err := cache.SetUPState(form.Uid, form.ProblemId, int(pb.SubmitState_UPStateExited), 60*2*time.Second); err != nil {
			logrus.Errorln(err.Error())
		}
	}()

	// 设置状态
	if err := cache.SetUPState(form.Uid, form.ProblemId, int(pb.SubmitState_UPStateNormal), 60*2*time.Second); err != nil {
		logrus.Errorln(err.Error())
		return
	}
	ok, param := preAction(form)
	if !ok {
		logrus.Errorln("预处理失败")
		return
	}

	start := time.Now()
	res := doAction(param)
	duration := time.Now().Sub(start).Milliseconds()
	logrus.Infof("---judge.Handle--- uid:%d, problemID:%d, total-cost:%d ms\n", form.Uid, form.ProblemId, duration)

	// 解锁用户
	//cache.UnLockUser(form.Uid)

	if res != nil {
		analyzeResult(param, res)

		// 记录结果
		data, err := json.Marshal(res)
		if err != nil {
			logrus.Errorln(err.Error())
			return
		}
		saveResult(param, data)
	}
}

func analyzeResult(param *types.Param, results []*pb.JudgeResult) {
	param.Ac = true
	for _, res := range results {
		if res.Status != "Accepted" {
			param.Ac = false
			break
		}
	}
}

func saveResult(param *types.Param, data []byte) {
	// 保存本次提交结果 2min过期
	err := cache.SetJudgeResult(param.Uid, param.ProblemID, data, 60*2*time.Second)
	if err != nil {
		logrus.Errorln(err.Error())
	}

	dbConn, err := registry.NewDBConnection(settings.Instance().RegistryConfig)
	if err != nil {
		logrus.Errorf("db服连接失败:%s\n", err.Error())
		return
	}
	defer dbConn.Close()

	// 更新用户提交记录
	err = updateUserSubmitRecord(param, data, dbConn)
	if err != nil {
		logrus.Errorln(err.Error())
	}
	// 更新用户AC记录
	err = updateUserAcProblemData(param, dbConn)
	if err != nil {
		logrus.Errorln(err.Error())
	}
	// 更新用户解题统计数据
	err = updateUserDoProblemStatistics(param, dbConn)
	if err != nil {
		logrus.Errorln(err.Error())
	}
}

func updateUserSubmitRecord(param *types.Param, data []byte, conn *grpc.ClientConn) error {
	client := pb.NewRecordServiceClient(conn)
	request := &pb.SaveUserSubmitRecordRequest{
		UserId:    param.Uid,
		ProblemId: param.ProblemID,
		Code:      param.Code,
		Result:    data,
		Lang:      param.Language,
		Stamp:     time.Now().Unix(),
	}

	_, err := client.SaveUserSubmitRecord(context.Background(), request)
	return err
}
func updateUserAcProblemData(param *types.Param, conn *grpc.ClientConn) error {
	client := pb.NewUserServiceClient(conn)
	request := &pb.UpdateUserACDataRequest{
		Uid:       param.Uid,
		ProblemId: param.ProblemID,
	}

	_, err := client.UpdateUserAcProblemData(context.Background(), request)
	return err
}
func updateUserDoProblemStatistics(param *types.Param, conn *grpc.ClientConn) error {
	client := pb.NewUserServiceClient(conn)
	request := &pb.UpdateUserDoProblemStatisticsRequest{
		Uid:             param.Uid,
		SubmitCountIncr: 1,
	}

	if param.Ac {
		request.AcCountIncr = 1
		switch param.Level {
		case 1:
			request.EasyCountIncr = 1
		case 2:
			request.MediumCountIncr = 1
		case 3:
			request.HardCountIncr = 1
		default:
			logrus.Infof("未知的题目难度:%d", param.Level)
		}
	}
	_, err := client.UpdateUserDoProblemStatistics(context.Background(), request)
	return err
}

func loadProblemDetail(problemID int64) (*pb.Problem, error) {
	dbConn, err := registry.NewDBConnection(settings.Instance().RegistryConfig)
	if err != nil {
		logrus.Errorf("db服连接失败:%s\n", err.Error())
		return nil, err
	}
	defer dbConn.Close()

	client := pb.NewProblemServiceClient(dbConn)
	response, err := client.GetProblemData(context.Background(), &pb.GetProblemRequest{Id: problemID})
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

func preAction(form *pb.SubmitForm) (bool, *types.Param) {
	param := &types.Param{}

	// 读取题目配置
	problemConfig, err := cache.GetProblemConfig(form.ProblemId)
	if err != nil {
		logrus.Errorln(err.Error())
		return false, nil
	}

	param.Uid = form.Uid
	param.ProblemID = form.ProblemId
	param.Code = form.Code
	param.Language = form.Lang
	param.ProblemConfig = problemConfig
	param.Level = problemConfig.Level
	return true, param
}

// 操作(编译，运行，评判)，操作的上下文信息需要缓存到redis
// redis需要持久化的信息：
// 1.本次提交的状态
// 2.编译结果
// 3.运行结果
// 4.评判结果
func doAction(param *types.Param) []*pb.JudgeResult {
	handler := NewHandler(settings.Instance().SandBox.Host, settings.Instance().SandBox.Port)

	results := make([]*pb.JudgeResult, 0)
	// 设置题目状态[编译]
	if err := cache.SetUPState(param.Uid, param.ProblemID, int(pb.SubmitState_UPStateCompiling), 60*2*time.Second); err != nil {
		logrus.Errorln(err.Error())
	}
	compileResult, err := handler.compile(param)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil
	}
	logrus.Debugln("编译结果:", compileResult)

	if compileResult.Status != "Accepted" {
		compileResult.Content = "编译失败"
		results = append(results, compileResult)
		// 更新状态
		if err := cache.SetUPState(param.Uid, param.ProblemID, int(pb.SubmitState_UPStateExited), 60*2*time.Second); err != nil {
			logrus.Errorln(err.Error())
			return nil
		}
		return results
	}
	compileResult.Content = "编译成功"
	results = append(results, compileResult)

	// 保存可执行文件的文件ID
	param.FileIds = compileResult.FileIds
	// 设置题目状态[判题中]
	if err := cache.SetUPState(param.Uid, param.ProblemID, int(pb.SubmitState_UPStateJudging), 60*2*time.Second); err != nil {
		logrus.Errorln(err.Error())
	}
	wgRun := new(sync.WaitGroup)
	wgRun.Add(1)
	goroutinePool.Instance().Submit(func() {
		defer wgRun.Done()
		handler.run(param)
	})
	wgJudge := new(sync.WaitGroup)
	wgJudge.Add(1)
	goroutinePool.Instance().Submit(func() {
		defer wgJudge.Done()
		judgeResults := handler.judge()
		results = append(results, judgeResults...)
	})

	wgRun.Wait()
	// 关闭管道，触发后续goroutine结束
	close(handler.runResults)
	wgJudge.Wait()

	return results
}
