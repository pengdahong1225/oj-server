package judgeService

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"question-service/models"
	"question-service/services/ants"
	"question-service/services/redis"
	"question-service/settings"
	"sync"
	"time"
)

var (
	baseUrl    string
	exeName    string      // 执行文件名
	runResults chan Result // 运行结果
)

type Param struct {
	uid           int64
	problemID     int64
	compileConfig *ProblemConfig
	runConfig     *ProblemConfig
	content       string // 源代码
	testCases     []TestCase
	fileIds       map[string]string // 文件id
}

func Init() {
	srv, err := settings.GetSystemConf("judge-service")
	if err != nil {
		logrus.Fatalln("Error getting system config:", err)
		return
	}
	baseUrl = fmt.Sprintf("http://%s:%d", srv.Host, srv.Port)

	runResults = make(chan Result, 256)
	exeName = "main"
}

// Handle 判题服务入口
func Handle(uid int64, form *models.SubmitForm) []Result {
	// 退出之后，需要将状态置为UPStateExited，主要针对异常退出的情况，正常退出会设置状态
	defer func() {
		if err := redis.SetUPState(uid, form.ProblemID, UPStateExited); err != nil {
			logrus.Errorln(err.Error())
		}
	}()

	// 设置“提交”状态
	if err := redis.SetUPState(uid, form.ProblemID, UPStateNormal); err != nil {
		logrus.Errorln(err.Error())
		return nil
	}

	ok, param := preAction(uid, form)
	if !ok {
		logrus.Errorln("预处理失败")
		return nil
	}

	start := time.Now()
	res := doAction(param)
	duration := time.Now().Sub(start).Milliseconds()
	logrus.Infoln("---judgeService.Handle--- uid:%d, problemID:%d, total-cost:%d ms", uid, form.ProblemID, duration)
	return res
}

func preAction(uid int64, form *models.SubmitForm) (bool, *Param) {
	param := &Param{}

	// 读取题目配置
	data, err := redis.GetProblemHotData(form.ProblemID)
	if err != nil {
		logrus.Errorln(err.Error())
		return false, nil
	}
	hotData := &ProblemHotData{}
	if err := json.Unmarshal([]byte(data), hotData); err != nil {
		logrus.Errorln(err.Error())
		return false, nil
	}
	compileConfig := &ProblemConfig{}
	if err := json.Unmarshal([]byte(hotData.CompileConfig), compileConfig); err != nil {
		logrus.Errorln(err.Error())
		return false, nil
	}
	runConfig := &ProblemConfig{}
	if err := json.Unmarshal([]byte(hotData.RunConfig), runConfig); err != nil {
		logrus.Errorln(err.Error())
		return false, nil
	}
	var testCases []TestCase
	if err := json.Unmarshal([]byte(hotData.TestCase), &testCases); err != nil {
		logrus.Errorln(err.Error())
		return false, nil
	}

	param.uid = uid
	param.problemID = form.ProblemID
	param.compileConfig = compileConfig
	param.runConfig = runConfig
	param.content = form.Code
	param.testCases = testCases
	return true, param
}

// 操作(编译，运行，评判)，操作的上下文信息需要缓存到redis
// redis需要持久化的信息：
// 1.本次提交的状态
// 2.编译结果
// 3.运行结果
// 4.评判结果
func doAction(param *Param) []Result {
	handler := &Handler{}
	results := make([]Result, 0)
	// 设置题目状态[编译]
	if err := redis.SetUPState(param.uid, param.problemID, UPStateCompiling); err != nil {
		logrus.Errorln(err.Error())
	}
	compileResult, err := handler.compile(param)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil
	}
	if compileResult.Status != "Accepted" {
		compileResult.Content = "编译失败"
		results = append(results, *compileResult)
		// 更新状态
		if err := redis.SetUPState(param.uid, param.problemID, UPStateExited); err != nil {
			logrus.Errorln(err.Error())
			return nil
		}
		return results
	}
	compileResult.Content = "编译成功"
	results = append(results, *compileResult)

	// 编译成功 => 设置题目状态[判题中]
	// 运行和判题是协同的步骤，由两个协程去同时进行，通过channel通信
	param.fileIds = compileResult.FileIds
	if err := redis.SetUPState(param.uid, param.problemID, UPStateJudging); err != nil {
		logrus.Errorln(err.Error())
	}
	wgRun := new(sync.WaitGroup)
	wgRun.Add(1)
	ants.AntsPoolInstance.Submit(func() {
		defer wgRun.Done()
		handler.run(param)
	})
	wgJudge := new(sync.WaitGroup)
	wgJudge.Add(1)
	ants.AntsPoolInstance.Submit(func() {
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
