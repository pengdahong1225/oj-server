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

type AdminHandler struct {
}

func (receiver AdminHandler) HandleAddQuestion(uid int64, form *models.AddProblemForm) *models.Response {
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
	request := &pb.CreateProblemRequest{Data: &pb.Problem{
		Title:       form.Title,
		Level:       form.Level,
		Tags:        form.Tags,
		Description: form.Desc,
		CreateBy:    uid,
		CompileConfig: &pb.ProblemConfig{
			ClockLimit:  form.CompileConfig.ClockLimit,
			CpuLimit:    form.CompileConfig.CpuLimit,
			MemoryLimit: form.CompileConfig.MemoryLimit,
			ProcLimit:   form.CompileConfig.ProcLimit,
		},
		RunConfig: &pb.ProblemConfig{
			ClockLimit:  form.RunConfig.ClockLimit,
			CpuLimit:    form.RunConfig.CpuLimit,
			MemoryLimit: form.RunConfig.MemoryLimit,
			ProcLimit:   form.RunConfig.ProcLimit,
		},
	}}
	for _, test := range form.TestCases {
		request.Data.TestCases = append(request.Data.TestCases, &pb.TestCase{
			Input:  test.Input,
			Output: test.Output,
		})
	}
	response, err := client.CreateProblemData(context.Background(), request)
	if err != nil {
		res.Code = http.StatusOK
		res.Message = "创建题目失败"
		logrus.Debugf("创建题目失败:%s\n", err.Error())
		return res
	}
	res.Code = http.StatusOK
	res.Message = "OK"
	res.Data = response.Id
	return res
}
