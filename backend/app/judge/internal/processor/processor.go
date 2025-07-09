package processor

import (
	"oj-server/app/judge/internal/define"
	"oj-server/proto/pb"
	"fmt"
	"oj-server/app/judge/internal/respository/cache"
	"github.com/sirupsen/logrus"
	"sync"
	"oj-server/module/gPool"
	"oj-server/consts"
	"time"
	"strings"
)

// 定义判题任务的处理模板
type IProcessor interface {
	Compile(param *define.Param) (*define.SandBoxApiResponse, error)
	Run(param *define.Param)
	Judge() []*pb.PBResult
}

// 基础实现
type BaseProcessor struct {
	impl       IProcessor // 注入具体实现
	param      *define.Param
	sandBoxUrl string
	runResults chan *define.RunResultInChan // 运行结果
}

func (b *BaseProcessor) Compile(param *define.Param) (*define.SandBoxApiResponse, error) {
	return nil, fmt.Errorf("not implemented")
}
func (b *BaseProcessor) Run(param *define.Param) {
	panic("run not implemented")
}
func (b *BaseProcessor) Judge() []*pb.PBResult {
	panic("judge not implemented")
}
func (b *BaseProcessor) CheckAnswer(X string, Y string) bool {
	X = strings.Replace(X, " ", "", -1)
	X = strings.Replace(X, "\n", "", -1)
	return X == Y
}

// HandleJudgeTask 模板方法
// 判题逻辑入口
func (b *BaseProcessor) HandleJudgeTask(form *pb.SubmitForm, param *define.Param) []*pb.PBResult {
	taskId := fmt.Sprintf("%d-%d", form.Uid, form.ProblemId)
	results := make([]*pb.PBResult, 0, 10)
	// 退出之后，需要将本次任务的状态置为UPStateExited，并且释放锁
	defer func() {
		if err := cache.SetTaskState(taskId, int(pb.SubmitState_UPStateExited), consts.TaskStateExpired); err != nil {
			logrus.Errorln(err.Error())
		}
		_ = cache.UnLockUser(form.Uid) // 释放锁
	}()
	// 初始化任务状态
	if err := cache.SetTaskState(taskId, int(pb.SubmitState_UPStateNormal), consts.TaskStateExpired); err != nil {
		logrus.Errorf("初始化任务状态失败, err=%s", err.Error())
		results = append(results, &pb.PBResult{
			Status: define.InternalError,
			ErrMsg: "初始化任务状态失败",
		})
		return results
	}
	// 设置题目状态[编译]
	err := cache.SetTaskState(taskId, int(pb.SubmitState_UPStateCompiling), consts.TaskStateExpired)
	if err != nil {
		logrus.Errorf("设置题目[%d]状态失败, err=%s", param.ProblemData.Id, err.Error())
		results = append(results, &pb.PBResult{
			Status: define.InternalError,
			ErrMsg: "task状态设置失败",
		})
		return results
	}
	compileResult, err := b.Compile(param)
	if err != nil {
		logrus.Errorf("请求编译发生意外, err=%s", err.Error())
		results = append(results, &pb.PBResult{
			Status: define.InternalError,
			ErrMsg: "请求编译发生意外",
		})
		return results
	}
	logrus.Debugln("编译结果:", *compileResult)

	// 需要将result类型转换为pb.PBResult
	pbResult := translatePBResult(compileResult)
	if compileResult.Status != define.Accepted {
		pbResult.Content = "Compile Error"
		results = append(results, pbResult)
		// 更新状态
		if err = cache.SetTaskState(taskId, int(pb.SubmitState_UPStateExited), consts.TaskStateExpired); err != nil {
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
	if err = cache.SetTaskState(taskId, int(pb.SubmitState_UPStateJudging), 60*2*time.Second); err != nil {
		logrus.Errorln(err.Error())
	}
	wgRun := new(sync.WaitGroup)
	wgRun.Add(1)
	_ = gPool.Instance().Submit(func() {
		defer wgRun.Done()
		b.Run(param)
	})
	wgJudge := new(sync.WaitGroup)
	wgJudge.Add(1)
	_ = gPool.Instance().Submit(func() {
		defer wgJudge.Done()
		judgeResults := b.Judge()
		results = append(results, judgeResults...)
	})

	wgRun.Wait()
	// 关闭管道，触发后续goroutine结束
	close(b.runResults)
	wgJudge.Wait()

	return results
}

func translatePBResult(resp *define.SandBoxApiResponse) *pb.PBResult {
	return &pb.PBResult{
		Status:     resp.Status,
		Content:    resp.ErrMsg,
		Memory:     resp.Memory,
		RunTime:    resp.RunTime,
		Time:       resp.Time,
		ExitStatus: resp.ExitStatus,
	}
}
