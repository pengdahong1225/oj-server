package judgeService

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"question-service/models"
	"question-service/services/redis"
	"question-service/settings"
	"time"
)

var (
	baseUrl   string
	exeName   string      // 执行文件名
	runResult chan Result // 运行结果
)

type Param struct {
	uid           int64
	problemID     int64
	compileConfig *CompileConfig
	runConfig     *RunConfig
	content       string // 源代码
	testCases     []TestCase
	fileIds       map[string]string // 文件id
}

func init() {
	srv, err := settings.GetSystemConf("judge-service")
	if err != nil {
		logrus.Fatalln("Error getting system config:", err)
		return
	}
	baseUrl = fmt.Sprintf("http://%s:%d", srv.Host, srv.Port)

	runResult = make(chan Result, 256)
	exeName = "main"
}

// Handle 判题服务入口
func Handle(uid int64, form *models.SubmitForm) {
	defer func() {
		// 恢复题目状态
		if err := redis.SetUPState(uid, form.ProblemID, UPStateNormal); err != nil {
			logrus.Errorln(err.Error())
		}
	}()

	// preAction
	// 1.题目缓存到redis，设置状态
	// 2.读取题目配置
	ok, param := preAction(uid, form)
	if !ok {
		logrus.Errorln("预处理失败!")
		return
	}

	start := time.Now()
	doAction(param)
	duration := time.Now().Sub(start).Milliseconds()
	logrus.Infoln("---judgeService.Handle--- uid:%d, problemID:%d, total-cost:%d ms", uid, form.ProblemID, duration)
}

func preAction(uid int64, form *models.SubmitForm) (bool, *Param) {
	param := &Param{}
	// 设置题目状态
	if err := redis.SetUPState(uid, form.ProblemID, UPStateNormal); err != nil {
		logrus.Errorln(err.Error())
		return false, nil
	}

	// 编译配置
	compileJson, err := redis.GetProblemCompileConfig(form.ProblemID)
	if err != nil {
		logrus.Errorln(err.Error())
		return false, nil
	}
	compileConfig := &CompileConfig{}
	if err := json.Unmarshal([]byte(compileJson), compileConfig); err != nil {
		logrus.Errorln(err.Error())
		return false, nil
	}

	// 运行配置
	runJson, err := redis.GetProblemRunConfig(form.ProblemID)
	if err != nil {
		logrus.Errorln(err.Error())
		return false, nil
	}
	runConfig := &RunConfig{}
	if err := json.Unmarshal([]byte(runJson), runConfig); err != nil {
		logrus.Errorln(err.Error())
		return false, nil
	}

	// 测试用例
	test, err := redis.GetProblemTest(form.ProblemID)
	if err != nil {
		logrus.Errorln(err.Error())
		return false, nil
	}
	var testCases []TestCase
	if err := json.Unmarshal([]byte(test), &testCases); err != nil {
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
func doAction(param *Param) {
	handler := &Handler{}
	// 设置题目状态[编译]
	if err := redis.SetUPState(param.uid, param.problemID, UPStateCompiling); err != nil {
		logrus.Errorln(err.Error())
	}
	compileResult, err := handler.Compile(param)
	if err != nil {
		logrus.Errorln(err.Error())
		return
	}
	// 判断编译结果，并更新状态和结果
	if compileResult.Status == "Accepted" {
		res := ResultInCache{
			Content: "编译成功",
			Results: []*Result{
				compileResult,
			},
		}
		bys, err := json.Marshal(&res)
		if err != nil {
			logrus.Errorln(err.Error())
			return
		}
		if err := redis.SetUPResult(param.uid, param.problemID, string(bys)); err != nil {
			logrus.Errorln(err.Error())
			return
		}
	} else {
		res := ResultInCache{
			Content: "编译失败",
			Results: []*Result{
				compileResult,
			},
		}
		bys, err := json.Marshal(&res)
		if err != nil {
			logrus.Errorln(err.Error())
			return
		}
		if err := redis.SetUPResult(param.uid, param.problemID, string(bys)); err != nil {
			logrus.Errorln(err.Error())
		}
		return
	}

	// 设置题目状态[运行]
	param.fileIds = compileResult.FileIds
	if err := redis.SetUPState(param.uid, param.problemID, UPStateRunning); err != nil {
		logrus.Errorln(err.Error())
	}
	handler.Run(param)

	// 设置题目状态[判题]
	if err := redis.SetUPState(param.uid, param.problemID, UPStateJudging); err != nil {
		logrus.Errorln(err.Error())
	}
	res, err := handler.Judge(param)
	if err != nil {
		logrus.Errorln(err.Error())
		return
	}
}
