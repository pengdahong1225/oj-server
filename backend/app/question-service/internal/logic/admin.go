package logic

import (
	"context"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/models"
	"github.com/pengdahong1225/oj-server/backend/module/registry"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
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
		Config: &pb.ProblemConfig{
			TestCases: nil,
			CompileLimit: &pb.Limit{
				CpuLimit:    form.Config.CompileLimit.CpuLimit,
				ClockLimit:  form.Config.CompileLimit.ClockLimit,
				MemoryLimit: form.Config.CompileLimit.MemoryLimit,
				StackLimit:  form.Config.CompileLimit.StackLimit,
				ProcLimit:   form.Config.CompileLimit.ProcLimit,
			},
			RunLimit: &pb.Limit{
				CpuLimit:    form.Config.RunLimit.CpuLimit,
				ClockLimit:  form.Config.RunLimit.ClockLimit,
				MemoryLimit: form.Config.RunLimit.MemoryLimit,
				StackLimit:  form.Config.RunLimit.StackLimit,
				ProcLimit:   form.Config.RunLimit.ProcLimit,
			},
		},
	}}
	for _, test := range form.Config.TestCases {
		request.Data.Config.TestCases = append(request.Data.Config.TestCases, &pb.TestCase{
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
