package judge

import (
	"context"
	"encoding/json"
	"github.com/pengdahong1225/oj-server/backend/app/judge-service/internal/cache"
	"github.com/pengdahong1225/oj-server/backend/app/judge-service/internal/types"
	"github.com/pengdahong1225/oj-server/backend/module/goroutinePool"
	"github.com/pengdahong1225/oj-server/backend/module/registry"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
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
	results := doAction(param)
	duration := time.Now().Sub(start).Milliseconds()
	logrus.Infof("---judge.Handle--- uid:%d, problemID:%d, total-cost:%d ms\n", form.Uid, form.ProblemId, duration)

	// 解锁用户
	cache.UnLockUser(form.Uid)

	if results != nil {
		analyzeResult(param, results)

		// 记录结果
		data, err := json.Marshal(results)
		if err != nil {
			logrus.Errorln(err.Error())
			return
		}
		saveResult(param, data)
	} else {
		logrus.Errorf("results is nil")
	}
}

// 分析结果
func analyzeResult(param *types.Param, results []*pb.PBResult) {
	param.Accepted = true
	param.Message = types.Accepted
	for _, result := range results {
		if result.Status != types.Accepted {
			param.Accepted = false
			param.Message = result.Content
			break
		}
	}
}

func saveResult(param *types.Param, data []byte) {
	// 保存本次提交结果 1min过期
	err := cache.SetJudgeResult(param.Uid, param.ProblemID, param.Message, 60*1*time.Second)
	if err != nil {
		logrus.Errorln(err.Error())
	}

	// 更新用户提交记录
	err = updateUserSubmitRecord(param, data)
	if err != nil {
		logrus.Errorln(err.Error())
	}
}

func updateUserSubmitRecord(param *types.Param, data []byte) error {
	conn, err := registry.NewDBConnection()
	if err != nil {
		logrus.Errorf("db服连接失败:%s\n", err.Error())
		return err
	}
	defer conn.Close()
	client := pb.NewRecordServiceClient(conn)
	request := &pb.SaveUserSubmitRecordRequest{
		Data: &pb.UserSubmitRecord{
			Uid:          param.Uid,
			ProblemId:    param.ProblemID,
			ProblemName:  param.ProblemConfig.Name,
			Status:       param.Message,
			Result:       data,
			Lang:         param.Language,
			ProblemLevel: param.ProblemConfig.Level,
			Code:         param.Code,
		},
	}

	_, err = client.SaveUserSubmitRecord(context.Background(), request)
	return err
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
	return true, param
}

// 操作(编译，运行，评判)，操作的上下文信息需要缓存到redis
// redis需要持久化的信息：
// 1.本次提交的状态
// 2.编译结果
// 3.运行结果
// 4.评判结果
func doAction(param *types.Param) []*pb.PBResult {
	handler := NewHandler(settings.Instance().SandBox.Host, settings.Instance().SandBox.Port)

	results := make([]*pb.PBResult, 0)
	// 设置题目状态[编译]
	if err := cache.SetUPState(param.Uid, param.ProblemID, int(pb.SubmitState_UPStateCompiling), 60*2*time.Second); err != nil {
		logrus.Errorln(err.Error())
	}
	compileResult, err := handler.compile(param)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil
	}
	logrus.Debugln("编译结果:", *compileResult)

	// 需要将result类型转换为pb.PBResult
	pbResult := translatePBResult(compileResult)
	if compileResult.Status != types.Accepted {
		pbResult.Content = "Compile Error"
		results = append(results, pbResult)
		// 更新状态
		if err := cache.SetUPState(param.Uid, param.ProblemID, int(pb.SubmitState_UPStateExited), 60*2*time.Second); err != nil {
			logrus.Errorln(err.Error())
			return nil
		}
		return results
	}
	pbResult.Content = "Compile Success"
	results = append(results, pbResult)

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

func translatePBResult(resp *types.SandBoxApiResponse) *pb.PBResult {
	return &pb.PBResult{
		Status:     resp.Status,
		Content:    resp.ErrMsg,
		Memory:     resp.Memory,
		RunTime:    resp.RunTime,
		Time:       resp.Time,
		ExitStatus: resp.ExitStatus,
	}
}
