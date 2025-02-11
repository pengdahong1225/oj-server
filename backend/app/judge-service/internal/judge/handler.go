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
	runResults chan *types.RunResultInChan // 运行结果
}

func NewHandler(host string, port int) *Handler {
	return &Handler{
		baseUrl:    fmt.Sprintf("http://%s:%d", host, port),
		runResults: make(chan *types.RunResultInChan, 100),
	}
}
func (r *Handler) compile(param *types.Param) (*types.SandBoxApiResponse, error) {
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
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorln("Error reading response:", err)
		return nil, err
	}
	logrus.Debugln("Response Status:", resp.Status)
	logrus.Debugln("Response Body:", string(respBody))

	var result []*types.SandBoxApiResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		logrus.Errorln(err.Error())
		return nil, err
	}
	if len(result) > 0 {
		// 编译结果就取第一条
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
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				logrus.Errorln("Error reading response:", err)
				return
			}
			logrus.Debugln("Response Status:", resp.Status)
			logrus.Debugln("Response Body:", string(respBody))

			// 将结果放入管道
			var results []*types.SandBoxApiResponse
			if err := json.Unmarshal(respBody, &results); err != nil {
				logrus.Errorln("Error unmarshalling JSON:", err)
				return
			}
			if len(results) > 0 {
				r.runResults <- &types.RunResultInChan{
					Result: results[0],
					Case:   test,
				}
			}
		})
	}
	wg.Wait()
}

// 检查结果状态，只check结果状态为Accepted的
func (r *Handler) judge() []*pb.PBResult {
	var results []*pb.PBResult
	for runResult := range r.runResults {
		pbResult := translatePBResult(runResult.Result)
		// status不为Accepted的，不用检测结果
		if runResult.Result.Status != types.Accepted {
			pbResult.Content = "Run Error"
			results = append(results, pbResult)
			continue
		}
		// 判断output是否满足预期
		// 不满足结果的状态为Wrong Answer
		if !r.checkAnswer(runResult.Result.Files["stdout"], runResult.Case.Output) {
			pbResult.Status = types.WrongAnswer
			pbResult.Content = "答案错误"
			results = append(results, pbResult)
		} else {
			pbResult.Status = types.Accepted
			pbResult.Content = "通过"
			results = append(results, pbResult)
		}
	}
	return results
}

func (r *Handler) checkAnswer(X string, Y string) bool {
	X = strings.Replace(X, " ", "", -1)
	X = strings.Replace(X, "\n", "", -1)
	return X == Y
}
