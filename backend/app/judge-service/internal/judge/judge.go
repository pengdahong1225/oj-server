package judge

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pengdahong1225/oj-server/backend/app/judge-service/internal/svc/cache"
	"github.com/pengdahong1225/oj-server/backend/app/judge-service/internal/types"
	"github.com/pengdahong1225/oj-server/backend/module/goroutinePool"
	"github.com/pengdahong1225/oj-server/backend/module/registry"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	baseUrl    string
	runResults chan *pb.JudgeResult // 运行结果
	once       sync.Once
)

func Init() {
	once.Do(func() {
		sanbox := settings.Instance().SandBox
		baseUrl = fmt.Sprintf("http://%s:%d", sanbox.Host, sanbox.Port)
		runResults = make(chan *pb.JudgeResult, 256)
	})
}

// Handle 判题服务入口
func Handle(form *pb.SubmitForm) {
	Init()

	// 退出之后，需要将"提交"状态置为UPStateExited
	defer func() {
		if err := cache.SetUPState(form.Uid, form.ProblemId, int(pb.SubmitState_UPStateExited)); err != nil {
			logrus.Errorln(err.Error())
		}
	}()

	// 设置“提交”状态
	if err := cache.SetUPState(form.Uid, form.ProblemId, int(pb.SubmitState_UPStateNormal)); err != nil {
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

	// 记录结果
	data, err := json.Marshal(res)
	if err != nil {
		logrus.Errorln(err.Error())
		return
	}
	saveResult(param, data)
}

func saveResult(param *types.Param, data []byte) {
	// 保存本次提交结果 2min过期
	err := cache.SetJudgeResult(param.Uid, param.ProblemID, data, 60*2*time.Second)
	if err != nil {
		logrus.Errorln(err.Error())
	}

	// 保存提交记录 record
	dbConn, err := registry.NewDBConnection(settings.Instance().RegistryConfig)
	if err != nil {
		logrus.Errorf("db服连接失败:%s\n", err.Error())
		return
	}
	defer dbConn.Close()

	client := pb.NewRecordServiceClient(dbConn)
	request := &pb.SaveUserSubmitRecordRequest{
		UserId:    param.Uid,
		ProblemId: param.ProblemID,
		Code:      param.Code,
		Result:    data,
		Lang:      param.Language,
		Stamp:     time.Now().Unix(),
	}

	_, err = client.SaveUserSubmitRecord(context.Background(), request)
	if err != nil {
		logrus.Errorln(err.Error())
	}
}

func preAction(form *pb.SubmitForm) (bool, *types.Param) {
	param := &types.Param{}

	// 读取题目配置
	problemConfig, err := cache.GetProblemConfig(form.ProblemId)
	if err != nil {
		logrus.Errorln("预处理失败:%s", err.Error())
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
func doAction(param *types.Param) []*pb.JudgeResult {
	handler := &Handler{}
	results := make([]*pb.JudgeResult, 0)
	// 设置题目状态[编译]
	if err := cache.SetUPState(param.Uid, param.ProblemID, int(pb.SubmitState_UPStateCompiling)); err != nil {
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
		if err := cache.SetUPState(param.Uid, param.ProblemID, int(pb.SubmitState_UPStateExited)); err != nil {
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
	if err := cache.SetUPState(param.Uid, param.ProblemID, int(pb.SubmitState_UPStateJudging)); err != nil {
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
	close(runResults)
	wgJudge.Wait()

	return results
}
