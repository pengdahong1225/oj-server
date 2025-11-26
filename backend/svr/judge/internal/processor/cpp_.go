package processor

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"oj-server/pkg/gPool"
	"oj-server/pkg/proto/pb"
	"oj-server/svr/judge/internal/biz"
	"sync"
)

type CPPProcessor struct {
	*BasicProcessor
}

func (cp *CPPProcessor) Compile(param *biz.Param) (*biz.SandBoxApiResponse, error) {
	form := biz.SandBoxApiForm{
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
	body := biz.SandBoxApiBody{}
	body.Cmd = append(body.Cmd, form)
	data, err := json.Marshal(body)
	if err != nil {
		logrus.Errorln("Error marshalling JSON:", err.Error())
		return nil, err
	}
	// POST请求
	req, err := http.NewRequest("POST", cp.sandBoxUrl+"/run", bytes.NewBuffer(data))
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
	logrus.Debugln("Response SandBoxApiBody:", string(respBody))

	var result []*biz.SandBoxApiResponse
	if err = json.Unmarshal(respBody, &result); err != nil {
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
func (cp *CPPProcessor) Run(param *biz.Param) {
	// 循环调用(并发地发送多个测试用例的运行请求，并等待所有请求完成)
	wg := new(sync.WaitGroup)

	for _, test := range param.ProblemConfig.TestCases {
		wg.Add(1)
		_ = gPool.Instance().Submit(func() {
			defer wg.Done()

			form := biz.SandBoxApiForm{
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

			body := biz.SandBoxApiBody{}
			body.Cmd = append(body.Cmd, form)
			data, err := json.Marshal(body)
			if err != nil {
				logrus.Errorln("Error marshalling JSON:", err)
				return
			}
			// POST请求
			client := &http.Client{}
			req, err := http.NewRequest("POST", cp.sandBoxUrl+"/run", bytes.NewBuffer(data))
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
			logrus.Debugln("Response SandBoxApiBody:", string(respBody))

			// 将结果放入管道
			var results []*biz.SandBoxApiResponse
			if err = json.Unmarshal(respBody, &results); err != nil {
				logrus.Errorln("Error unmarshalling JSON:", err)
				return
			}
			if len(results) > 0 {
				cp.runResultsChan <- &biz.RunResultInChan{
					Result: results[0],
					Case:   test,
				}
			}
		})
	}
	wg.Wait()
}

// 检查结果状态，只check结果状态为Accepted的
func (cp *CPPProcessor) Judge() []*pb.JudgeResultItem {
	var results []*pb.JudgeResultItem
	for runResult := range cp.runResultsChan {
		pbResult := translatePBResult(runResult.Result)
		// status不为Accepted的，不用检测结果
		if runResult.Result.Status != biz.Accepted {
			pbResult.Content = "Run Error"
			results = append(results, pbResult)
			continue
		}
		// 判断output是否满足预期
		// 不满足结果的状态为Wrong Answer
		if !cp.CheckAnswer(runResult.Result.Files["stdout"], runResult.Case.Output) {
			pbResult.Status = biz.WrongAnswer
			pbResult.Content = "答案错误"
			results = append(results, pbResult)
		} else {
			pbResult.Status = biz.Accepted
			pbResult.Content = "通过"
			results = append(results, pbResult)
		}
	}
	return results
}
