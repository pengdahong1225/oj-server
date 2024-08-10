package handler

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"question-service/internal/proto"
	"question-service/models"
	"question-service/services/ants"
	"question-service/services/judgeService"
	"question-service/services/redis"
	"question-service/services/registry"
	"question-service/settings"
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
	if state != judgeService.UserStateNormal {
		res.Code = http.StatusBadRequest
		res.Message = "用户处于判题状态，请稍等..."
		return res
	}

	// 设置用户状态为判题中
	if err := redis.SetUserState(uid, judgeService.UserStateJudging); err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		logrus.Errorln(err.Error())
		return res
	}

	// 异步处理：提交到judgeService
	err = ants.AntsPoolInstance.Submit(func() {
		// 处理结果
		results := judgeService.Handle(uid, form)
		data, err := json.Marshal(results)
		if err != nil {
			logrus.Errorln(err.Error())
		} else {
			// 存储结果
			if err := redis.SetJudgeResult(uid, form.ProblemID, string(data)); err != nil {
				logrus.Errorln(err.Error())
			}
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
	if state != judgeService.UPStateExited {
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
	var results []judgeService.Result
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