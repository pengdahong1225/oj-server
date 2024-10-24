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
)

type ProblemLogic struct {
}

func (receiver ProblemLogic) GetProblemList(params *models.QueryProblemListParams) *models.Response {
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

	res.Message = "OK"
	res.Data = response
	return res
}

// HandleProblemSubmit
// 判断“用户”是否处于判题状态？true就拒绝
// 用户提交了题目就立刻返回，并给题目设置状态
// 客户端通过其他接口轮询题目结果
func (receiver ProblemLogic) HandleProblemSubmit(uid int64, form *models.SubmitForm) *models.Response {
	res := &models.Response{
		Code:    models.Success,
		Message: "",
		Data:    nil,
	}

	ok, err := cache.LockUser(uid, 60)
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
		consts.RabbitMqCommentQueue,
		consts.RabbitMqCommentKey,
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

func (receiver ProblemLogic) HandleProblemDetail(problemID int64) *models.Response {
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

func (receiver ProblemLogic) HandleQueryResult(uid int64, problemID int64) *models.Response {
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
		res.Message = "running"
		return res
	}

	// 查询结果
	r, err := cache.QueryJudgeResult(uid, problemID)
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorln(err.Error())
		return res
	}
	// 解析
	var results []models.SubmitResult
	if err := json.Unmarshal([]byte(r), &results); err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		logrus.Errorln(err.Error())
		return res
	}

	res.Message = "OK"
	res.Data = results
	return res
}
