package processor

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"oj-server/global"
	"oj-server/pkg/gPool"
	"oj-server/proto/pb"
	"oj-server/svr/gateway/internal/configs"
	"oj-server/svr/judge/internal/biz"
	"strings"
	"sync"
)

// 定义判题任务的处理模板
type IProcessor interface {
	Compile(param *biz.Param) (*biz.SandBoxApiResponse, error)
	Run(param *biz.Param)
	Judge() []*pb.PBResult
}

// 处理器工厂
func NewProcessor(language string, uc *biz.JudgeUseCase) (*BasicProcessor, error) {
	var processor IProcessor

	language = strings.ToLower(language)
	switch language {
	case "c":
		processor = &CProcessor{}
	case "cpp":
		processor = &CProcessor{}
	case "go":
		processor = &GoProcessor{}
	case "python":
		processor = &PyProcessor{}
	default:
		return nil, fmt.Errorf("language not supported, language=%s", language)
	}

	// 查询sandbox地址
	var addr string
	for _, item := range configs.AppConf.SandBoxCfg {
		if item.Type == language {
			addr = fmt.Sprintf("http://%s:%d", item.Host, item.Port)
			break
		}
	}
	if addr == "" {
		return nil, fmt.Errorf("target sandbox config not found, language=%s", language)
	}

	return &BasicProcessor{
		uc:             uc,
		impl:           processor,
		sandBoxUrl:     addr,
		runResultsChan: make(chan *biz.RunResultInChan, 100),
	}, nil
}

// 模板处理器
type BasicProcessor struct {
	impl           IProcessor                // 注入具体实现
	param          *biz.Param                // 上下文参数
	sandBoxUrl     string                    // 判题沙箱地址
	uc             *biz.JudgeUseCase         // 仓库
	runResultsChan chan *biz.RunResultInChan // 判题结果
}

// 模板方法
// 判题逻辑入口
func (b *BasicProcessor) HandleJudgeTask(form *pb.SubmitForm, param *biz.Param) []*pb.PBResult {
	taskId := fmt.Sprintf("%d-%d", form.Uid, form.ProblemId)
	results := make([]*pb.PBResult, 0, 10)
	// 退出之后，需要将本次任务的状态置为UPStateExited，并且释放锁
	defer func() {
		_ = b.uc.SetTaskState(taskId, int(pb.SubmitState_UPStateExited))
		key := fmt.Sprintf("%s:%d", global.UserLockPrefix, form.Uid)
		_ = b.uc.UnLock(key) // 释放锁
	}()
	// 初始化任务状态
	err := b.uc.SetTaskState(taskId, int(pb.SubmitState_UPStateNormal))
	if err != nil {
		logrus.Errorf("初始化任务状态失败, err=%s", err.Error())
		results = append(results, &pb.PBResult{
			Status: biz.InternalError,
			ErrMsg: "初始化任务状态失败",
		})
		return results
	}
	// 设置题目状态[编译]
	err = b.uc.SetTaskState(taskId, int(pb.SubmitState_UPStateCompiling))
	if err != nil {
		logrus.Errorf("设置题目[%d]状态失败, err=%s", param.ProblemData.Id, err.Error())
		results = append(results, &pb.PBResult{
			Status: biz.InternalError,
			ErrMsg: "task状态设置失败",
		})
		return results
	}
	compileResult, err := b.impl.Compile(param)
	if err != nil {
		logrus.Errorf("请求编译发生意外, err=%s", err.Error())
		results = append(results, &pb.PBResult{
			Status: biz.InternalError,
			ErrMsg: "请求编译发生意外",
		})
		return results
	}
	logrus.Debugln("编译结果:", *compileResult)

	// 需要将result类型转换为pb.PBResult
	pbResult := translatePBResult(compileResult)
	if compileResult.Status != biz.Accepted {
		pbResult.Content = "Compile Error"
		results = append(results, pbResult)
		// 更新状态
		err = b.uc.SetTaskState(taskId, int(pb.SubmitState_UPStateExited))
		if err != nil {
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
	// 运行和判题并行处理
	err = b.uc.SetTaskState(taskId, int(pb.SubmitState_UPStateJudging))
	if err != nil {
		logrus.Errorln(err.Error())
	}
	wgRun := new(sync.WaitGroup)
	wgRun.Add(1)
	_ = gPool.Instance().Submit(func() {
		defer wgRun.Done()
		b.impl.Run(param)
	})
	wgJudge := new(sync.WaitGroup)
	wgJudge.Add(1)
	_ = gPool.Instance().Submit(func() {
		defer wgJudge.Done()
		judgeResults := b.impl.Judge()
		results = append(results, judgeResults...)
	})

	wgRun.Wait()
	// 关闭管道，触发后续goroutine结束
	close(b.runResultsChan)
	wgJudge.Wait()

	return results
}
func (b *BasicProcessor) CheckAnswer(X string, Y string) bool {
	X = strings.Replace(X, " ", "", -1)
	X = strings.Replace(X, "\n", "", -1)
	return X == Y
}
func translatePBResult(resp *biz.SandBoxApiResponse) *pb.PBResult {
	return &pb.PBResult{
		Status:     resp.Status,
		Content:    resp.ErrMsg,
		Memory:     resp.Memory,
		RunTime:    resp.RunTime,
		Time:       resp.Time,
		ExitStatus: resp.ExitStatus,
	}
}
