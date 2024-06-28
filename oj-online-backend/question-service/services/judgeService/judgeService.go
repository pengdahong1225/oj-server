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
	baseUrl       string
	compileResult []chan string // 编译结果
	runResult     []chan string // 运行结果
	judgeResult   []chan string // 判题结果
)

type Param struct {
	compileConfig *CompileConfig
	runConfig     *RunConfig
	test          string
}

func init() {
	srv, err := settings.GetSystemConf("judge-service")
	if err != nil {
		logrus.Fatalln("Error getting system config:", err)
		return
	}
	baseUrl = fmt.Sprintf("http://%s:%d", srv.Host, srv.Port)

	compileResult = make([]chan string, 1024)
	runResult = make([]chan string, 1024)
	judgeResult = make([]chan string, 1024)
}

// Handle 判题服务入口
func Handle(uid int64, form *models.SubmitForm) {
	// preAction
	// 1.题目缓存到redis，设置状态
	// 2.读取题目配置
	ok, param := preAction(uid, form)
	if !ok {
		logrus.Errorln("预处理失败!")
		return
	}

	start := time.Now()

	judger := &Handler{}

	// 编译协程
	wgCompile := new(sync.WaitGroup)
	wgCompile.Add(1)
	ants.AntsPoolInstance.Submit(func() {
		defer wgCompile.Done()
		judger.Compile(param)
	})

	// 运行协程
	wgRun := new(sync.WaitGroup)
	wgRun.Add(1)
	ants.AntsPoolInstance.Submit(func() {
		defer wgRun.Done()
		judger.Run(param)
	})

	// 评判协程
	wgJudge := new(sync.WaitGroup)
	wgJudge.Add(1)
	ants.AntsPoolInstance.Submit(func() {
		defer wgJudge.Done()
		judger.Judge(param)
	})

	wgCompile.Wait()
	wgRun.Wait()
	wgJudge.Wait()

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
	// 读取题目配置
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
	compileConfig.Content = form.Code

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

	test, err := redis.GetProblemTest(form.ProblemID)
	if err != nil {
		logrus.Errorln(err.Error())
		return false, nil
	}

	param.compileConfig = compileConfig
	param.runConfig = runConfig
	param.test = test
	return true, param
}
