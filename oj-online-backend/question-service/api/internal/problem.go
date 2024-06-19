package internal

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	pb "question-service/api/proto"
	"question-service/models"
	"question-service/services/judgeClient"
	"question-service/services/registry"
	"question-service/settings"
)

func ProblemSet(cursor int) *models.Response {
	res := &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	}
	dbConn, err := registry.NewDBConnection(settings.Conf.RegistryConfig)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorf("db服连接失败:%s\n", err.Error())
		return res
	}
	defer dbConn.Close()

	client := pb.NewDBServiceClient(dbConn)
	request := &pb.GetProblemListRequest{Cursor: int32(cursor)}
	response, err := client.GetProblemList(context.Background(), request)
	if err != nil {
		res.Code = http.StatusOK
		res.Message = "获取题目列表失败"
		logrus.Debugf("获取题目列表失败:%s\n", err.Error())
		return res
	}

	res.Code = http.StatusOK
	res.Message = "OK"
	res.Data = response
	return res
}

func ProblemSubmitHandler(uid int64, form *models.SubmitForm) *models.Response {
	res := &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	}

	system, err := settings.GetSystemConf("judge-service")
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorln(err.Error())
		return res
	}

	// 从缓存读取test_cast
	// test, err := redis.GetTestCaseJson(uid)
	// if err != nil {
	// 	res.Code = http.StatusInternalServerError
	// 	res.Message = err.Error()
	// 	logrus.Errorln(err.Error())
	// 	return res
	// }
	test := "\"test_case\": \"{\\r\\n  \\\"info\\\": {\\r\\n    \\\"test_case_number\\\": 1,\\r\\n    \\\"spj\\\": false,\\r\\n    \\\"test_cases\\\": {\\r\\n      \\\"1\\\": {\\r\\n        \\\"input_name\\\": \\\"1.in\\\",\\r\\n        \\\"output_name\\\": \\\"1.out\\\"\\r\\n      }\\r\\n    }\\r\\n  },\\r\\n  \\\"input\\\": [\\r\\n    {\\r\\n      \\\"name\\\": \\\"1.in\\\",\\r\\n      \\\"content\\\": \\\"1 2\\\"\\r\\n    }\\r\\n  ],\\r\\n  \\\"output\\\": [\\r\\n    {\\r\\n      \\\"name\\\": \\\"1.out\\\",\\r\\n      \\\"content\\\": \\\"3\\\"\\r\\n    }\\r\\n  ]\\r\\n}\",\n"

	request := &pb.SSJudgeRequest{
		Code:         form.Code,
		Language:     form.Lang,
		TestCaseJson: test,
	}
	client := judgeClient.TcpClient{}

	if err := client.Connect(fmt.Sprintf("%s:%d", system.Host, system.Port)); err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorln(err.Error())
		return res
	}
	response, err := client.RpcJudgeRequest(request)
	if err != nil {
		res.Code = http.StatusOK
		res.Message = err.Error()
		logrus.Errorf("RpcJudgeRequest err:%s", err.Error())
		return res
	}

	res.Data = response.ResultList
	return res
}

func GetProblemDetailHandler(problemID int64) *models.Response {
	res := &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	}
	dbConn, err := registry.NewDBConnection(settings.Conf.RegistryConfig)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorf("db服连接失败:%s\n", err.Error())
		return res
	}
	defer dbConn.Close()

	client := pb.NewDBServiceClient(dbConn)
	response, err := client.GetProblemData(context.Background(), &pb.GetProblemRequest{
		Id: problemID,
	})
	if err != nil {
		res.Code = http.StatusOK
		res.Message = err.Error()
		logrus.Errorln(err.Error())
		return res
	}

	res.Code = http.StatusOK
	res.Message = "OK"
	res.Data = response.Data

	return res
}
