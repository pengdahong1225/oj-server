package handler

import (
	"context"
	"encoding/json"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/models"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/services/goroutinePool"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/services/mq"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/services/redis"
	"github.com/pengdahong1225/Oj-Online-Server/app/question-service/settings"
	"github.com/pengdahong1225/Oj-Online-Server/pkg/registry"
	pb "github.com/pengdahong1225/Oj-Online-Server/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"net/http"
)

type ProblemHandler struct {
}

func (receiver ProblemHandler) HandleProblemSet(cursor int) *models.Response {
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

// HandleProblemSubmit
// 判断“用户”是否处于判题状态？true就拒绝
// 用户提交了题目就立刻返回，并给题目设置状态
// 客户端通过其他接口轮询题目结果
func (receiver ProblemHandler) HandleProblemSubmit(uid int64, form *models.SubmitForm) *models.Response {
	res := &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	}

	// 判断用户是否处于判题状态
	state, err := redis.GetUserState(uid)
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
	if err := redis.SetUserState(uid, int(pb.UserState_UserStateJudging)); err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorln(err.Error())
		return res
	}

	// 异步处理
	err = goroutinePool.PoolInstance.Submit(func() {
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
		} else {
			// 提交到mq
			productor := &mq.Producer{
				Exkind:     "direct",
				Exname:     "judge",
				QuName:     "judge",
				RoutingKey: "judge",
			}
			productor.Publish(data)
		}
	})

	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorln(err.Error())
		return res
	} else {
		// 返回题目id
		res.Code = http.StatusOK
		res.Message = "题目提交成功"
		res.Data = map[string]interface{}{
			"problemID": form.ProblemID,
		}
		return res
	}
}

func (receiver ProblemHandler) HandleProblemDetail(problemID int64) *models.Response {
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

// 先查询状态：如果没有查询到，就意味着最近没有提交题目or题目提交过期了
// 如果是已退出状态，就可以查询结果
// 如果是：有状态，但是没有结果 -> 被中断，需要查看日志排查
func (receiver ProblemHandler) HandleQueryResult(uid int64, problemID int64) *models.Response {
	res := &models.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	}

	// 查询状态
	state, err := redis.QueryUPState(uid, problemID)
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
	r, err := redis.QueryJudgeResult(uid, problemID)
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

func (receiver ProblemHandler) HandleProblemSearch(name string) *models.Response {
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
