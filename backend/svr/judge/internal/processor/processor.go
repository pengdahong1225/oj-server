package processor

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"oj-server/pkg/gPool"
	"oj-server/pkg/proto/pb"
	"oj-server/svr/judge/internal/biz"
	"oj-server/svr/judge/internal/configs"
	"strings"
	"sync"
)

// 定义判题任务的处理模板
type IProcessor interface {
	Compile(param *biz.Param) (*biz.SandBoxApiResponse, error)
	Run(param *biz.Param)
	Judge() []*pb.JudgeResultItem
}

// 处理器工厂
func NewProcessor(language string) (*BasicProcessor, error) {
	basicProcessor := &BasicProcessor{
		runResultsChan: make(chan *biz.RunResultInChan, 100),
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
	logrus.Debugf("sandbox url: %s", addr)
	basicProcessor.sandBoxUrl = addr

	// 生成处理器
	var processor IProcessor
	language = strings.ToLower(language)
	switch language {
	case "c":
		processor = &CProcessor{
			BasicProcessor: basicProcessor,
		}
	case "cpp":
		processor = &CPPProcessor{
			BasicProcessor: basicProcessor,
		}
	case "go":
		processor = &GoProcessor{
			BasicProcessor: basicProcessor,
		}
	case "python":
		processor = &PyProcessor{
			BasicProcessor: basicProcessor,
		}
	default:
		return nil, fmt.Errorf("language not supported, language=%s", language)
	}
	basicProcessor.impl = processor

	return basicProcessor, nil
}

// 模板处理器
type BasicProcessor struct {
	impl           IProcessor                // 注入具体实现
	param          *biz.Param                // 上下文参数
	sandBoxUrl     string                    // 判题沙箱地址
	runResultsChan chan *biz.RunResultInChan // 判题结果
}

// 模板方法
// 判题逻辑入口
func (b *BasicProcessor) HandleJudgeTask(task *pb.JudgeSubmission, param *biz.Param) []*pb.JudgeResultItem {
	results := make([]*pb.JudgeResultItem, 0, 10)

	// 编译
	compileResult, err := b.impl.Compile(param)
	if err != nil {
		logrus.Errorf("请求编译发生意外, err=%s", err.Error())
		results = append(results, &pb.JudgeResultItem{
			Status: biz.InternalError,
			ErrMsg: "编译发生意外",
		})
		return results
	}

	// 需要将result类型转换为pb.PBResult
	result := translatePBResult(compileResult)
	if compileResult.Status != biz.Accepted {
		result.Content = "Compile Error"
		results = append(results, result)
		return results
	}
	result.Content = "Compile Success"
	results = append(results, result)

	// 保存可执行文件的文件ID
	param.FileIds = compileResult.FileIds

	// 判题
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
func translatePBResult(resp *biz.SandBoxApiResponse) *pb.JudgeResultItem {
	return &pb.JudgeResultItem{
		Status:     resp.Status,
		Content:    resp.ErrMsg,
		Memory:     resp.Memory,
		RunTime:    resp.RunTime,
		Time:       resp.Time,
		ExitStatus: resp.ExitStatus,
	}
}
