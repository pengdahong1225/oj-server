package logic

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	pb "question-service/logic/proto"
	"question-service/models"
	"question-service/services/judgeClient"
	"question-service/services/redis"
	"question-service/services/registry"
	"question-service/settings"
	"question-service/utils"
	"question-service/views"
)

func QuestionSet(ctx *gin.Context, cursor int32) {
	dbConn, err := registry.NewDBConnection(settings.Conf.RegistryConfig)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "db服务连接失败",
		})
	}
	defer dbConn.Close()

	client := pb.NewDBServiceClient(dbConn)
	request := &pb.GetQuestionListRequest{Cursor: cursor}
	response, err := client.GetQuestionList(context.Background(), request)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	var questionList []views.Question
	for _, v := range response.Data {
		questionList = append(questionList, views.Question{
			Id:    v.Id,
			Title: v.Title,
			Level: v.Level,
			Tags:  v.Tags,
		})
	}
	data, _ := json.Marshal(questionList)
	ctx.JSON(http.StatusOK, gin.H{
		"data":   data,
		"cursor": response.Cursor,
	})
}

func QuestionDetail(ctx *gin.Context, id int64) {
	dbConn, err := registry.NewDBConnection(settings.Conf.RegistryConfig)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "db服务连接失败",
		})
	}
	defer dbConn.Close()

	client := pb.NewDBServiceClient(dbConn)
	request := &pb.GetQuestionRequest{Id: id}
	response, err := client.GetQuestionData(context.Background(), request)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	data := views.Question{
		Id:          response.Data.Id,
		Title:       response.Data.Title,
		Level:       response.Data.Level,
		Tags:        response.Data.Tags,
		Description: response.Data.Description,
		TestCase:    response.Data.TestCase,
		Template:    response.Data.Template,
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func QuestionQuery(ctx *gin.Context, name string) {
	dbConn, err := registry.NewDBConnection(settings.Conf.RegistryConfig)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "db服务连接失败",
		})
	}
	defer dbConn.Close()

	client := pb.NewDBServiceClient(dbConn)
	request := &pb.QueryQuestionWithNameRequest{Name: name}

	response, err := client.QueryQuestionWithName(context.Background(), request)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	if response == nil {
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}
	var questionList []views.Question
	for _, v := range response.Data {
		questionList = append(questionList, views.Question{
			Id:          v.Id,
			Title:       v.Title,
			Level:       v.Level,
			Tags:        v.Tags,
			Description: v.Description,
			TestCase:    v.TestCase,
			Template:    v.Template,
		})
	}
	data, _ := json.Marshal(questionList)
	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func QuestionRun(ctx *gin.Context, form *models.QuestionForm) {
	// 获取server端dsn
	dsn, err := registry.GetJudgeServerDsn(settings.Conf.RegistryConfig)
	if err != nil {
		logrus.Infoln("获取judge-server端dsn失败:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "获取judge-server端dsn失败",
			"err": err.Error(),
		})
		return
	}

	// todo 获取测试用例
	test_case_json, err := redis.GetTestCaseJson(form.Id)
	if err != nil {
		logrus.Infoln("获取测试用例失败:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "获取测试用例失败",
			"err": err.Error(),
		})
		return
	}
	// 生成提交id
	submit_id, err := utils.GenerateSubmitID(int(form.UserId), int(form.Id))
	if err != nil {
		logrus.Infoln("生成提交id失败:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成提交id失败",
			"err": err.Error(),
		})
		return
	}

	// todo 请求judge-service
	request := &pb.SSJudgeRequest{
		Code:         form.Code,
		Language:     form.Clang,
		TestCaseJson: test_case_json,
		SubmitId:     submit_id,
	}
	client := judgeClient.TcpClient{}
	err = client.Connect(dsn)
	if err != nil {
		logrus.Infoln("连接Judge服务失败:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "连接Judge服务失败",
			"err": err.Error(),
		})
		return
	}
	result, err := client.Request(request)
	if err != nil {
		logrus.Errorln("发送消息失败:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "发送消息失败",
			"err": err.Error(),
		})
		return
	}

	// TODO 返回结果
	response := &views.QuestionResponse{
		QuestionID: form.Id,
		UserID:     form.UserId,
		Clang:      form.Clang,
	}
	for _, value := range result.ResultList {
		response.ResultList = append(response.ResultList, views.QuestionResult{
			Result:   value.Result,
			RealTime: value.RealTime,
			CpuTime:  value.CpuTime,
			Memory:   value.Memory,
			Signal:   value.Signal,
			ExitCode: value.ExitCode,
			Error:    value.Error,
			Content:  value.Content,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": response,
	})
}

func QuestionSubmit(ctx *gin.Context, form *models.QuestionForm) {

}

// 保存提交记录
func updateSubmitRecord(msg []byte, result string) {
	dbConn, err := registry.NewDBConnection(settings.Conf.RegistryConfig)
	if err != nil {
		logrus.Errorf("db服务连接失败:%s", err.Error())
		return
	}
	defer dbConn.Close()

	// 解析msg
	var question models.QuestionForm
	if err := json.Unmarshal(msg, &question); err != nil {
		logrus.Errorf("解析msg失败:%s", err.Error())
		return
	}

	client := pb.NewDBServiceClient(dbConn)
	request := &pb.UpdateUserSubmitRecordRequest{
		UserId:     question.UserId,
		QuestionId: question.Id,
		Code:       question.Code,
		Lang:       question.Clang,
		Result:     result,
	}

	_, err = client.UpdateUserSubmitRecord(context.Background(), request)
	if err != nil {
		logrus.Errorf("更新提交记录失败:%s", err.Error())
	}
}
