package logic

import (
	"context"
	"encoding/json"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/models"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/svc/cache"
	"github.com/pengdahong1225/oj-server/backend/consts"
	"github.com/pengdahong1225/oj-server/backend/module/mq"
	"github.com/pengdahong1225/oj-server/backend/module/registry"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"slices"
	"time"
)

type ProblemLogic struct {
}

func (r ProblemLogic) UpdateQuestion(uid int64, form *models.UpdateProblemForm, config *models.ProblemConfig) *models.Response {
	res := &models.Response{
		Code:    models.Success,
		Message: "",
		Data:    nil,
	}
	dbConn, err := registry.NewDBConnection(settings.Instance().RegistryConfig)
	if err != nil {
		res.Code = models.Failed
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
				CpuLimit:    config.CompileLimit.CpuLimit,
				ClockLimit:  config.CompileLimit.ClockLimit,
				MemoryLimit: config.CompileLimit.MemoryLimit,
				StackLimit:  config.CompileLimit.StackLimit,
				ProcLimit:   config.CompileLimit.ProcLimit,
			},
			RunLimit: &pb.Limit{
				CpuLimit:    config.RunLimit.CpuLimit,
				ClockLimit:  config.RunLimit.ClockLimit,
				MemoryLimit: config.RunLimit.MemoryLimit,
				StackLimit:  config.RunLimit.StackLimit,
				ProcLimit:   config.RunLimit.ProcLimit,
			},
		},
	}}
	for _, test := range config.TestCases {
		request.Data.Config.TestCases = append(request.Data.Config.TestCases, &pb.TestCase{
			Input:  test.Input,
			Output: test.Output,
		})
	}
	response, err := client.UpdateProblemData(context.Background(), request)
	if err != nil {
		res.Code = models.Failed
		res.Message = "update题目失败"
		logrus.Debugf("update题目失败:%s\n", err.Error())
		return res
	}
	res.Code = models.Success
	res.Message = "OK"
	res.Data = response.Id
	return res
}

func (r ProblemLogic) DeleteQuestion(problemID int64) *models.Response {
	res := &models.Response{
		Code:    models.Success,
		Message: "",
		Data:    nil,
	}
	dbConn, err := registry.NewDBConnection(settings.Instance().RegistryConfig)
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorf("db服连接失败:%s\n", err.Error())
		return res
	}
	defer dbConn.Close()
	client := pb.NewProblemServiceClient(dbConn)

	_, err = client.DeleteProblemData(context.Background(), &pb.DeleteProblemRequest{Id: problemID})
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		return res
	}
	res.Message = "OK"
	return res
}

func (r ProblemLogic) GetProblemList(params *models.QueryProblemListParams, uid int64) *models.Response {
	res := &models.Response{
		Code:    models.Success,
		Message: "",
		Data:    nil,
	}
	dbConn, err := registry.NewDBConnection(settings.Instance().RegistryConfig)
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorf("db服连接失败:%s\n", err.Error())
		return res
	}
	defer dbConn.Close()

	client := pb.NewProblemServiceClient(dbConn)
	request := &pb.GetProblemListRequest{
		Page:     params.Page,
		PageSize: params.PageSize,
		Keyword:  params.Keyword,
		Tag:      params.Tag,
	}
	response, err := client.GetProblemList(context.Background(), request)
	if err != nil {
		res.Code = models.Failed
		res.Message = "获取题目列表失败"
		logrus.Debugf("获取题目列表失败:%s\n", err.Error())
		return res
	}

	if uid > 0 {
		param := &models.UPSSParams{}
		param.Uid = uid
		for _, v := range response.Data {
			param.ProblemIds = append(param.ProblemIds, v.Id)
		}
		list, err := User{}.queryUserSolvedListByProblemList(param)
		if err != nil {
			logrus.Errorln(err.Error())
		} else {
			for _, v := range response.Data {
				if slices.Contains(list, v.Id) {
					v.Status = 1
				}
			}
		}
	}

	res.Message = "OK"
	res.Data = response
	return res
}

// OnProblemSubmit
// 判断“用户”是否处于判题状态？true就拒绝
// 用户提交了题目就立刻返回，并给题目设置状态
// 客户端通过其他接口轮询题目结果
func (r ProblemLogic) OnProblemSubmit(uid int64, form *models.SubmitForm) *models.Response {
	res := &models.Response{
		Code:    models.Success,
		Message: "",
		Data:    nil,
	}

	ok, err := cache.LockUser(uid, 60*time.Second)
	if err != nil {
		logrus.Errorln("lock user err:", err.Error())
		res.Code = models.Failed
		res.Message = err.Error()
		return res
	}
	if !ok {
		logrus.Errorln("lock user failed")
		res.Code = models.Failed
		res.Message = "正在判题中"
		return res
	}

	// 异步处理
	// protobuf序列化
	pbForm := pb.SubmitForm{
		Uid:       uid,
		ProblemId: form.ProblemID,
		Title:     form.Title,
		Lang:      form.Lang,
		Code:      form.Code,
	}
	data, err := proto.Marshal(&pbForm)
	if err != nil {
		logrus.Errorln(err.Error())
		res.Code = models.Failed
		res.Message = err.Error()
		cache.UnLockUser(uid)
		return res
	}
	// 提交到mq
	productor := mq.NewProducer(
		consts.RabbitMqExchangeKind,
		consts.RabbitMqExchangeName,
		consts.RabbitMqJudgeQueue,
		consts.RabbitMqJudgeKey,
	)
	if !productor.Publish(data) {
		res.Code = models.Failed
		res.Message = "任务提交mq失败"
		logrus.Errorln("任务提交mq失败")
		cache.UnLockUser(uid)
		return res
	} else {
		res.Message = "题目提交成功"
		res.Data = map[string]interface{}{
			"problemID": form.ProblemID,
		}
		return res
	}
}

func (r ProblemLogic) GetProblemDetail(problemID int64) *models.Response {
	res := &models.Response{
		Code:    models.Success,
		Message: "",
		Data:    nil,
	}
	dbConn, err := registry.NewDBConnection(settings.Instance().RegistryConfig)
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorf("db服连接失败:%s\n", err.Error())
		return res
	}
	defer dbConn.Close()

	client := pb.NewProblemServiceClient(dbConn)
	response, err := client.GetProblemData(context.Background(), &pb.GetProblemRequest{
		Id: problemID,
	})
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorln(err.Error())
		return res
	}

	res.Message = "OK"
	res.Data = response.Data

	return res
}

func (r ProblemLogic) QueryResult(uid int64, problemID int64) *models.Response {
	res := &models.Response{
		Code:    models.Success,
		Message: "",
		Data:    nil,
	}

	// 查询状态
	state, err := cache.QueryUPState(uid, problemID)
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorln(err.Error())
		return res
	}
	if state < 0 {
		res.Code = models.Failed
		res.Message = "题目未提交或提交已过期"
		return res
	}
	if state != int32(pb.SubmitState_UPStateExited) {
		res.Code = models.Failed
		res.Message = "running..."
		return res
	}

	// 当state == SubmitState_UPStateExited时，从缓存中提取结果
	bys, err := cache.GetJudgeResult(uid, problemID)
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorf("获取结果failed: %s\n", err.Error())
		return res
	}
	data := models.QuerySubmitResultResponse{
		Uid:       uid,
		ProblemID: problemID,
		Result:    nil,
	}
	err = json.Unmarshal([]byte(bys), &data.Result)
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorf("反序列化失败: %s\n", err.Error())
		return res
	}

	res.Message = "OK"
	res.Data = data
	return res
}
