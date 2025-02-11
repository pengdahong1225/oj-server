package judge

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pengdahong1225/oj-server/backend/app/judge-service/internal/types"
	"github.com/pengdahong1225/oj-server/backend/module/goroutinePool"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"sync"
)

type Handler struct {
	baseUrl    string
	runResults chan *pb.JudgeResult // 运行结果
}

func NewHandler(host string, port int) *Handler {
	return &Handler{
		baseUrl:    fmt.Sprintf("http://%s:%d", host, port),
		runResults: make(chan *pb.JudgeResult, 100),
	}
}
func (r *Handler) compile(param *types.Param) (*pb.JudgeResult, error) {
	form := types.SandBoxApiForm{
		CpuLimit:    param.ProblemConfig.CompileLimit.CpuLimit,
		ClockLimit:  param.ProblemConfig.CompileLimit.ClockLimit,
		MemoryLimit: param.ProblemConfig.CompileLimit.MemoryLimit,
		StackLimit:  param.ProblemConfig.CompileLimit.StackLimit,
		ProcLimit:   param.ProblemConfig.CompileLimit.ProcLimit,
		Args:        []string{"/usr/bin/g++", "main.cc", "-o", "main"},
		Env:         []string{"PATH=/usr/bin:/bin"},
	}
	form.Files = []map[string]any{
		{"content": ""},
		{"name": "stdout", "max": 10240},
		{"name": "stderr", "max": 10240},
	}
	form.CopyIn = map[string]map[string]string{
		"main.cc": {"content": param.Code},
	}
	form.CopyOut = []string{"stdout", "stderr"}
	form.CopyOutCached = []string{"main"}

	// 构造body
	body := types.Body{}
	body.Cmd = append(body.Cmd, form)
	data, err := json.Marshal(body)
	if err != nil {
		logrus.Errorln("Error marshalling JSON:", err.Error())
		return nil, err
	}
	// POST请求
	req, err := http.NewRequest("POST", r.baseUrl+"/run", bytes.NewBuffer(data))
	if err != nil {
		logrus.Errorln("Error creating request:", err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// 创建HTTP客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorln("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorln("Error reading response:", err)
		return nil, err
	}
	logrus.Debugln("Response Status:", resp.Status)
	logrus.Debugln("Response Body:", string(bodyResp))

	var result []*pb.JudgeResult
	if err := json.Unmarshal(bodyResp, &result); err != nil {
		logrus.Errorln(err.Error())
		return nil, err
	}
	if len(result) > 0 {
		return result[0], nil
	} else {
		logrus.Errorln("result len = 0")
		return nil, errors.New("result len = 0")
	}
}

func (r *Handler) run(param *types.Param) {
	// 循环调用(并发地发送多个测试用例的运行请求，并等待所有请求完成)
	wg := new(sync.WaitGroup)

	for _, test := range param.ProblemConfig.TestCases {
		wg.Add(1)
		goroutinePool.Instance().Submit(func() {
			defer wg.Done()

			form := types.SandBoxApiForm{
				CpuLimit:    param.ProblemConfig.CompileLimit.CpuLimit,
				ClockLimit:  param.ProblemConfig.CompileLimit.ClockLimit,
				MemoryLimit: param.ProblemConfig.CompileLimit.MemoryLimit,
				StackLimit:  param.ProblemConfig.CompileLimit.StackLimit,
				ProcLimit:   param.ProblemConfig.CompileLimit.ProcLimit,
				Args:        []string{"main"},
				Env:         []string{"PATH=/usr/bin:/bin"},
			}
			form.Files = []map[string]any{
				{"content": test.Input},
				{"name": "stdout", "max": 10240},
				{"name": "stderr", "max": 10240}}
			form.CopyIn = map[string]map[string]string{
				"main": {"fileId": param.FileIds["main"]},
			}

			body := types.Body{}
			body.Cmd = append(body.Cmd, form)
			data, err := json.Marshal(body)
			if err != nil {
				logrus.Errorln("Error marshalling JSON:", err)
				return
			}
			// POST请求
			client := &http.Client{}
			req, err := http.NewRequest("POST", r.baseUrl+"/run", bytes.NewBuffer(data))
			if err != nil {
				logrus.Errorln("Error creating request:", err)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				logrus.Errorln("Error sending request:", err)
				return
			}
			defer resp.Body.Close()

			// 读取响应
			bodyResp, err := io.ReadAll(resp.Body)
			if err != nil {
				logrus.Errorln("Error reading response:", err)
				return
			}
			logrus.Debugln("Response Status:", resp.Status)
			logrus.Debugln("Response Body:", string(bodyResp))

			// 将结果放入管道
			var result []*pb.JudgeResult
			if err := json.Unmarshal(bodyResp, &result); err != nil {
				logrus.Errorln("Error unmarshalling JSON:", err)
				return
			}

			if len(result) > 0 {
				result[0].TestCase = test
				r.runResults <- result[0]
			}
		})
	}
	wg.Wait()
}

// 检查结果状态，只check结果状态为Accepted的
func (r *Handler) judge() []*pb.JudgeResult {
	var results []*pb.JudgeResult
	for runResult := range r.runResults {
		if runResult.Status != "Accepted" {
			runResult.Content = runResult.Status
			results = append(results, runResult)
			continue
		}
		// 判断output是否满足预期
		// 不满足结果的状态为Wrong Answer
		if !r.checkAnswer(runResult.Files["stdout"], runResult.TestCase.Output) {
			runResult.Status = "Wrong Answer"
			runResult.Content = "Wrong Answer"
			results = append(results, runResult)
		} else {
			runResult.Status = "Accepted"
			runResult.Content = "Accepted"
			results = append(results, runResult)
		}
	}
	return results
}

func (r *Handler) checkAnswer(X string, Y string) bool {
	X = strings.Replace(X, " ", "", -1)
	X = strings.Replace(X, "\n", "", -1)
	return X == Y
}
