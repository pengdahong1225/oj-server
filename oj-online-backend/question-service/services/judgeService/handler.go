package judgeService

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Handler struct {
}

func (receiver *Handler) Compile(param *Param) (string, error) {
	return "", nil
}

func (receiver *Handler) Run(param *Param) {
	return
}

func (receiver *Handler) Judge(param *Param) (string, error) {
	return "", nil
}

func Compile(config *CompileConfig) (string, error) {
	// 定义请求的body内容
	body := map[string]interface{}{
		"cmd": []map[string]interface{}{
			{
				"args": []string{"/usr/bin/g++", "a.cc", "-o", "a"},
				"env":  []string{"PATH=/usr/bin:/bin"},
				"files": []map[string]interface{}{
					{"content": ""},
					{"name": "stdout", "max": 10240},
					{"name": "stderr", "max": 10240}},
				"cpuLimit":      config.CpuLimit,
				"memoryLimit":   config.MemoryLimit,
				"procLimit":     config.ProcLimit,
				"copyIn":        map[string]map[string]string{"a.cc": {"content": config.Content}},
				"copyOut":       []string{"stdout", "stderr"},
				"copyOutCached": []string{"a"},
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

func Run(config *RunConfig) {
	// 定义请求的body内容
	body := map[string]interface{}{
		"cmd": []map[string]interface{}{
			{
				"args": []string{"a"},
				"env":  []string{"PATH=/usr/bin:/bin"},
				"files": []map[string]interface{}{
					{"content": config.StdIn},
					{"name": "stdout", "max": 10240},
					{"name": "stderr", "max": 10240}},
				"cpuLimit":    config.CpuLimit,
				"memoryLimit": config.MemoryLimit,
				"procLimit":   config.ProcLimit,
				"copyIn":      map[string]map[string]string{"a": {"fileId": config.FileId}},
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
}
