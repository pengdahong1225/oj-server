package judge

import (
	"bytes"
	"encoding/json"
	"errors"
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
}

func (receiver *Handler) compile(param *types.Param) (*pb.JudgeResult, error) {
	// POST请求
	body := map[string]any{
		"cmd": []map[string]any{
			{
				// 资源限制
				"cpuLimit":    param.ProblemConfig.CompileLimit.CpuLimit,
				"memoryLimit": param.ProblemConfig.CompileLimit.MemoryLimit,
				"procLimit":   param.ProblemConfig.CompileLimit.ProcLimit,
				// 程序命令行参数
				"args": []string{"/usr/bin/g++", "main.cc", "-o", "main"},
				// 程序环境变量
				"env": []string{"PATH=/usr/bin:/bin"},
				// 指定 标准输入、标准输出和标准错误的文件 (null 是为了 pipe 的使用情况准备的，而且必须被 pipeMapping 的 in / out 指定)
				"files": []map[string]any{
					{"content": ""},
					{"name": "stdout", "max": 10240},
					{"name": "stderr", "max": 10240},
				},
				// 在执行程序之前复制进容器的文件列表
				"copyIn": map[string]map[string]string{"main.cc": {"content": param.Code}},
				// 在执行程序后从容器文件系统中复制出来的文件列表(不返回结果的内容，返回一个文件id)
				"copyOut":       []string{"stdout", "stderr"},
				"copyOutCached": []string{"main"},
			},
		},
	}
	data, err := json.Marshal(body)
	if err != nil {
		logrus.Errorln("Error marshalling JSON:", err.Error())
		return nil, err
	}
	req, err := http.NewRequest("POST", baseUrl+"/run", bytes.NewBuffer(data))
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

func (receiver *Handler) run(param *types.Param) {
	// 循环调用(协程池并发地发送多个请求，并等待所有请求完成)
	wg := new(sync.WaitGroup)

	for _, test := range param.ProblemConfig.TestCases {
		wg.Add(1)
		goroutinePool.Instance().Submit(func() {
			defer wg.Done()

			body := map[string]any{
				"cmd": []map[string]any{
					{
						// 资源限制
						"cpuLimit":    param.ProblemConfig.RunLimit.CpuLimit,
						"memoryLimit": param.ProblemConfig.RunLimit.MemoryLimit,
						"procLimit":   param.ProblemConfig.RunLimit.ProcLimit,
						// 程序命令行参数
						"args": []string{"main"},
						// 程序环境变量
						"env": []string{"PATH=/usr/bin:/bin"},
						// 指定 标准输入、标准输出和标准错误的文件 (null 是为了 pipe 的使用情况准备的，而且必须被 pipeMapping 的 in / out 指定)
						"files": []map[string]any{
							{"content": test.Input},
							{"name": "stdout", "max": 10240},
							{"name": "stderr", "max": 10240}},
						// 在执行程序之前复制进容器的文件列表（这个缓存文件的 ID 来自上一个请求返回的 fileIds）
						"copyIn": map[string]map[string]string{"main": {"fileId": param.FileIds["main"]}},
					},
				},
			}
			data, err := json.Marshal(body)
			if err != nil {
				logrus.Errorln("Error marshalling JSON:", err)
				return
			}
			// POST请求
			client := &http.Client{}
			req, err := http.NewRequest("POST", baseUrl+"/run", bytes.NewBuffer(data))
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
				runResults <- result[0]
			}
		})
	}
	wg.Wait()
}

// 1.检查结果状态，只judge结果状态为Accepted的
// 2.如果结果状态为其他，不judge直接缓存
func (receiver *Handler) judge() []*pb.JudgeResult {
	var results []*pb.JudgeResult
	for runResult := range runResults {
		if runResult.Status != "Accepted" {
			runResult.Content = "可执行程序运行错误"
			results = append(results, runResult)
			continue
		}
		// 判断output是否满足预期
		// 不满足结果的状态为Wrong Answer
		if !receiver.analyzeAnswer(runResult.Files["stdout"], runResult.TestCase.Output) {
			runResult.Status = "Wrong Answer"
			runResult.Content = "运行结果错误"
			results = append(results, runResult)
		} else {
			runResult.Content = "运行结果正确"
			results = append(results, runResult)
		}
	}
	return results
}

func (receiver *Handler) analyzeAnswer(X string, Y string) bool {
	X = strings.Replace(X, " ", "", -1)
	X = strings.Replace(X, "\n", "", -1)
	return X == Y
}
