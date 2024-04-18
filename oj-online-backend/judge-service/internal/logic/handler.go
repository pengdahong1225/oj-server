package logic

import (
	"encoding/json"
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"judge-service/internal/logic/impl"
	"judge-service/models"
	"judge-service/services/ants"
	"judge-service/services/redis"
	"os"
	"sync"
	"time"
)

// Handle 入口
func Handle(form *models.JudgeRequest) []byte {
	handler := NewHandler()
	rsp := handler.JudgeQuestion(form)
	msg, _ := json.Marshal(rsp)
	return msg
}

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (receiver *Handler) JudgeQuestion(form *models.JudgeRequest) *models.JudgeBack {
	rsp := &models.JudgeBack{
		SessionID:  form.SessionID,
		QuestionID: form.QuestionID,
		UserID:     form.UserID,
		Clang:      form.Clang,
		Status:     -1,
		Tips:       "",
		Output:     "",
	}

	// 获取测试用例
	var cases []models.TestCase
	caseData, err := getTestCase(form.QuestionID)
	if err != nil || caseData == "" {
		logrus.Errorf("error: %s", err.Error())
		rsp.Status = models.EN_Status_Internal
		rsp.Tips = "系统错误，没有找到测试用例"
		return rsp
	}
	if err = json.Unmarshal([]byte(caseData), &cases); err != nil {
		logrus.Errorf("error: %s", err.Error())
		rsp.Status = models.EN_Status_Internal
		rsp.Tips = "系统错误，没有找到测试用例"
		return rsp
	}

	// 新建沙箱
	sandBox, err := impl.NewSandBox(form.Clang)
	if err != nil {
		logrus.Errorf("error: %s", err.Error())
		rsp.Status = models.EN_Status_Internal
		rsp.Tips = "不支持当前语言"
		rsp.Output = err.Error()
		return rsp
	}

	// 检查代码合法性
	if err = sandBox.CheckCodeValid([]byte(form.Code)); err != nil {
		logrus.Errorf("error: %s", err.Error())
		rsp.Status = models.EN_Status_CompileError
		rsp.Tips = "代码非法"
		rsp.Output = err.Error()
		return rsp
	}

	// 代码保存本地
	path, err := saveCode(form.UserID, form.Clang, form.Code)
	if err != nil {
		logrus.Errorf("error: %s", err.Error())
		rsp.Status = models.EN_Status_Internal
		rsp.Tips = "保存代码失败"
		rsp.Output = err.Error()
		return rsp
	}

	// 运行代码
	var passCount int
	wg := new(sync.WaitGroup)
	wg.Add(1)
	err = ants.AntsPoolInstance.Submit(func() {
		defer wg.Done()
		for _, testCase := range cases {
			// 入口
			if !sandBox.Run(path, rsp, testCase) {
				break
			}
			passCount++
		}
	})
	if err != nil {
		logrus.Errorf("error: %s", err.Error())
		rsp.Status = models.EN_Status_Internal
		rsp.Tips = "运行代码失败"
		rsp.Output = err.Error()
		return rsp
	}
	timer := time.AfterFunc(time.Second*5, func() {
		wg.Done()
		rsp.Status = models.EN_Status_TimeOut
		rsp.Tips = "运行超时"
	})

	// 等待运行
	wg.Wait()
	timer.Stop()

	if passCount == len(cases) {
		rsp.Status = models.EN_Status_OK
		rsp.Tips = "通过"
		rsp.Output = caseData
	}
	logrus.Infof("uid:%d,qid:%d,passCount:%d", form.UserID, form.QuestionID, passCount)
	return rsp
}

func getTestCase(subKey int64) (string, error) {
	var key = "QuestionTestCases"
	// 获取测试用例
	conn := redis.NewConn()
	defer conn.Close()

	reply, err := redigo.String(conn.Do("HGET", key, subKey))
	if err != nil {
		return "", err
	}
	return reply, nil
}

// 保存代码
func saveCode(uid int64, lang string, code string) (string, error) {
	dirName := fmt.Sprintf("code/%d", uid)
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		return "", err
	}
	var path string
	switch lang {
	case "c":
		path = dirName + "/main.c"
	case "cpp":
		path = dirName + "/main.cpp"
	case "java":
		path = dirName + "/main.java"
	case "python":
		path = dirName + "/main.py"
	case "go":
		path = dirName + "/main.go"
	case "js":
		path = dirName + "/main.js"
	case "php":
		path = dirName + "/main.php"
	default:
		return "", fmt.Errorf("unsupported language: %s", lang)
	}
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = file.Write([]byte(code))
	if err != nil {
		return "", err
	}
	return path, nil
}
