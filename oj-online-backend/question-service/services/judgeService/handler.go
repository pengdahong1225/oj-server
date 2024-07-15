package judgeService

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"question-service/models"
	"question-service/services/ants"
	"question-service/services/redis"
)

type Handler struct {
}

func (receiver *Handler) Compile(param *Param) (string, error) {
	// 设置题目状态
	if err := redis.SetUPState(param.uid, param.problemID, models.UPStateCompiling); err != nil {
		logrus.Errorln(err.Error())
	}
	// 定义请求的body内容
	body := map[string]interface{}{
		"cmd": []map[string]interface{}{
			{
				// 程序命令行参数
				"args": []string{"/usr/bin/g++", "a.cc", "-o", "a"},
				// 程序环境变量
				"env": []string{"PATH=/usr/bin:/bin"},
				// 指定 标准输入、标准输出和标准错误的文件 (null 是为了 pipe 的使用情况准备的，而且必须被 pipeMapping 的 in / out 指定)
				"files": []map[string]interface{}{
					{"content": ""},
					{"name": "stdout", "max": 10240},
					{"name": "stderr", "max": 10240}},
				// 资源限制
				"cpuLimit":    param.compileConfig.CpuLimit,
				"memoryLimit": param.compileConfig.MemoryLimit,
				"procLimit":   param.compileConfig.ProcLimit,
				// 在执行程序之前复制进容器的文件列表
				"copyIn": map[string]map[string]string{"a.cc": {"content": param.content}},
				// 在执行程序后从容器文件系统中复制出来的文件列表(不返回结果的内容，返回一个文件id)
				"copyOut":       []string{"stdout", "stderr"},
				"copyOutCached": []string{exeName},
			},
		},
	}
	// 将body内容编码为JSON
	jsonData, err := json.Marshal(body)
	if err != nil {
		logrus.Errorln("Error marshalling JSON:", err.Error())
		return "", err

	}

	// 创建HTTP POST请求
	req, err := http.NewRequest("POST", baseUrl+"/run", bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.Errorln("Error creating request:", err.Error())
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 创建HTTP客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorln("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorln("Error reading response:", err)
		return "", err
	}

	// 打印响应
	logrus.Debugln("Response Status:", resp.Status)
	logrus.Debugln("Response Body:", string(bodyResp))

	return string(bodyResp), nil
}

func (receiver *Handler) Run(param *Param) (string, error) {
	// 循环调用(扔协程池一次性调用所有的请求，再一起等待)
	for _, test := range param.testCases {
		ants.AntsPoolInstance.Submit(func() {
			// 定义请求的body内容
			body := map[string]interface{}{
				"cmd": []map[string]interface{}{
					{
						// 程序命令行参数
						"args": []string{"a"},
						// 程序环境变量
						"env": []string{"PATH=/usr/bin:/bin"},
						// 指定 标准输入、标准输出和标准错误的文件 (null 是为了 pipe 的使用情况准备的，而且必须被 pipeMapping 的 in / out 指定)
						"files": []map[string]interface{}{
							{"content": test.Input},
							{"name": "stdout", "max": 10240},
							{"name": "stderr", "max": 10240}},
						// 资源限制
						"cpuLimit":    param.runConfig.CpuLimit,
						"memoryLimit": param.runConfig.MemoryLimit,
						"procLimit":   param.runConfig.ProcLimit,
						// 在执行程序之前复制进容器的文件列表（这个缓存文件的 ID 来自上一个请求返回的 fileIds）
						"copyIn": map[string]map[string]string{exeName: {"fileId": param.fileIds[exeName]}},
					},
				},
			}

			// 将body内容编码为JSON
			jsonData, err := json.Marshal(body)
			if err != nil {
				logrus.Errorln("Error marshalling JSON:", err)
				return
			}

			// 创建HTTP POST请求
			req, err := http.NewRequest("POST", baseUrl+"/run", bytes.NewBuffer(jsonData))
			if err != nil {
				logrus.Errorln("Error creating request:", err)
				return
			}

			// 设置请求头
			req.Header.Set("Content-Type", "application/json")

			// 创建HTTP客户端并发送请求
			client := &http.Client{}
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

			// 打印响应
			logrus.Debugln("Response Status:", resp.Status)
			logrus.Debugln("Response Body:", string(bodyResp))
		})
	}

	return "", nil
}

func (receiver *Handler) Judge(param *Param) (string, error) {
	return "", nil
}
