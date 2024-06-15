package internal

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	pb "question-service/api/proto"
	"question-service/models"
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

	return res
}
