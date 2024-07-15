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
	baseUrl       string
	compileResult chan string // 编译结果
	runResult     chan string // 运行结果
	judgeResult   chan string // 判题结果
	exeName       string      // 执行文件名
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

	compileResult = make(chan string, 1024)
	runResult = make(chan string, 1024)
	judgeResult = make(chan string, 1024)
	exeName = "main"
}

// Handle 判题服务入口
func Handle(uid int64, form *models.SubmitForm) {
	defer func() {
		// 恢复题目状态
		if err := redis.SetUPState(uid, form.ProblemID, models.UPStateNormal); err != nil {
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

	handler := &Handler{}
	// 编译
	res, err := handler.Compile(param)
	if err != nil {
		logrus.Errorln(err.Error())
		return
	}
	ok, compileResult := checkResult(res)
	if !ok {
		logrus.Errorln("result校验失败")
		return
	}
	// 运行
	param.fileIds = compileResult.FileIds

	res, err = handler.Run(param)
	if err != nil {
		logrus.Errorln(err.Error())
		return
	}
	// 判题
	res, err = handler.Judge(param)
	if err != nil {
		logrus.Errorln(err.Error())
		return
	}

	duration := time.Now().Sub(start).Milliseconds()
	logrus.Infoln("---judgeService.Handle--- uid:%d, problemID:%d, total-cost:%d ms", uid, form.ProblemID, duration)
}

func preAction(uid int64, form *models.SubmitForm) (bool, *Param) {
	param := &Param{}
	// 设置题目状态
	if err := redis.SetUPState(uid, form.ProblemID, models.UPStateNormal); err != nil {
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

// 结果判断
func checkResult(res string) (bool, *Result) {
	result := &Result{}
	if err := json.Unmarshal([]byte(res), result); err != nil {
		logrus.Errorln(err.Error())
		return false, nil
	}
	if result.Status != "Accepted" {
		return false, nil
	}
	return true, result
}
