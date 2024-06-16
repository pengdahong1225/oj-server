package internal

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	pb "question-service/api/proto"
	"question-service/models"
	"question-service/services/judgeClient"
	"question-service/services/redis"
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
	test, err := redis.GetTestCaseJson(uid)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorln(err.Error())
		return res
	}

	request := &pb.SSJudgeRequest{
		Code:         form.Code,
		Language:     form.Lang,
		TestCaseJson: test,
	}
	client := judgeClient.TcpClient{}

	if err := client.Connect(system.Host); err != nil {
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
