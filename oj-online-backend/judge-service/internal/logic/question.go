package logic

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"judge-service/global"
	"judge-service/internal/models"
	"os"
)

// QuestionRun 代码运行
func QuestionRun(form *models.QuestionForm) []byte {
	cases := getTestCase(form.Id)
	if cases == "" {
		return nil
	}
	// 代码保存本地
	codePath, err := saveCode(form.UserId, form.Clang, form.Code)
	if err != nil {
		return nil
	}
	// 检查代码

	// 运行代码

	return nil
}

func getTestCase(subKey int64) string {
	var key = "QuestionTestCases"
	// 获取测试用例
	conn := global.RedisPoolInstance.Get()
	defer conn.Close()

	reply, err := redis.String(conn.Do("HGET", key, subKey))
	if err != nil {
		logrus.Errorf("getTestCase error: %s", err.Error())
		return ""
	}
	return reply
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
