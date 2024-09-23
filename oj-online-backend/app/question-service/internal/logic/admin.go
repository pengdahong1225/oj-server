package logic

import (
	"context"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/models"
	"github.com/pengdahong1225/Oj-Online-Server/module/registry"
	"github.com/pengdahong1225/Oj-Online-Server/module/settings"
	"github.com/pengdahong1225/Oj-Online-Server/proto/pb"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AdminLogic struct {
}

func (receiver AdminLogic) HandleUpdateQuestion(uid int64, form *models.AddProblemForm) *models.Response {
	res := &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	}
	dbConn, err := registry.NewDBConnection(settings.Instance().RegistryConfig)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorf("db服连接失败:%s\n", err.Error())
		return res
	}
	defer dbConn.Close()

	client := pb.NewProblemServiceClient(dbConn)
	request := &pb.UpdateProblemRequest{Data: &pb.Problem{
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
	response, err := client.UpdateProblemData(context.Background(), request)
	if err != nil {
		res.Code = http.StatusOK
		res.Message = "update题目失败"
		logrus.Debugf("update题目失败:%s\n", err.Error())
		return res
	}
	res.Code = http.StatusOK
	res.Message = "OK"
	res.Data = response.Id
	return res
}

func (receiver AdminLogic) HandleDelQuestion(problemID int64) *models.Response {
	res := &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	}
	dbConn, err := registry.NewDBConnection(settings.Instance().RegistryConfig)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorf("db服连接失败:%s\n", err.Error())
		return res
	}
	defer dbConn.Close()
	client := pb.NewProblemServiceClient(dbConn)

	_, err = client.DeleteProblemData(context.Background(), &pb.DeleteProblemRequest{Id: problemID})
	if err != nil {
		res.Message = err.Error()
		return res
	}
	res.Message = "OK"
	return res
}
