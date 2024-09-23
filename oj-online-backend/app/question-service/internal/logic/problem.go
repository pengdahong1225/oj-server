package logic

import (
	"context"
	"encoding/json"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/models"
	redis2 "github.com/pengdahong1225/Oj-Online-Server/app/question-service/internal/svc/redis"
	"github.com/pengdahong1225/Oj-Online-Server/consts"
	"github.com/pengdahong1225/Oj-Online-Server/module/mq"
	"github.com/pengdahong1225/Oj-Online-Server/module/registry"
	"github.com/pengdahong1225/Oj-Online-Server/module/settings"
	"github.com/pengdahong1225/Oj-Online-Server/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"net/http"
)

type ProblemLogic struct {
}

func (receiver ProblemLogic) HandleProblemSet(cursor int) *models.Response {
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

// HandleProblemSubmit
// 判断“用户”是否处于判题状态？true就拒绝
// 用户提交了题目就立刻返回，并给题目设置状态
// 客户端通过其他接口轮询题目结果
func (receiver ProblemLogic) HandleProblemSubmit(uid int64, form *models.SubmitForm) *models.Response {
	res := &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	}

	// 判断用户是否处于判题状态
	state, err := redis2.GetUserState(uid)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorln(err.Error())
		return res
	}
	if state != int(pb.UserState_UserStateNormal) {
		res.Code = http.StatusBadRequest
		res.Message = "用户处于判题状态，请稍等..."
		return res
	}
	// 设置用户状态为判题中
	if err := redis2.SetUserState(uid, int(pb.UserState_UserStateJudging)); err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorln(err.Error())
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
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
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
		res.Code = http.StatusInternalServerError
		logrus.Errorln("任务提交mq失败")
		return res
	} else {
		res.Code = http.StatusOK
		res.Message = "题目提交成功"
		res.Data = map[string]interface{}{
			"problemID": form.ProblemID,
		}
		return res
	}
}

func (receiver ProblemLogic) HandleProblemDetail(problemID int64) *models.Response {
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

// 先查询状态：如果没有查询到，就意味着最近没有提交题目or题目提交过期了
// 如果是已退出状态，就可以查询结果
// 如果是：有状态，但是没有结果 -> 被中断，需要查看日志排查
func (receiver ProblemLogic) HandleQueryResult(uid int64, problemID int64) *models.Response {
	res := &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	}

	// 查询状态
	state, err := redis2.QueryUPState(uid, problemID)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorln(err.Error())
		return res
	}
	if state < 0 {
		res.Code = http.StatusOK
		res.Message = "题目未提交或提交已过期"
		return res
	}
	if state != int(pb.SubmitState_UPStateExited) {
		res.Code = http.StatusOK
		res.Message = "running"
		return res
	}

	// 查询结果
	r, err := redis2.QueryJudgeResult(uid, problemID)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorln(err.Error())
		return res
	}
	// 解析
	var results []models.SubmitResult
	if err := json.Unmarshal([]byte(r), &results); err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorln(err.Error())
		return res
	}

	res.Code = http.StatusOK
	res.Message = "OK"
	res.Data = results
	return res
}

func (receiver ProblemLogic) HandleProblemSearch(name string) *models.Response {
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
	response, err := client.QueryProblemWithName(context.Background(), &pb.QueryProblemWithNameRequest{Name: name})
	if err != nil {
		res.Code = http.StatusOK
		res.Message = err.Error()
		logrus.Errorln(err.Error())
		return res
	}
	res.Message = "OK"
	res.Data = response.Data
	return res
}
