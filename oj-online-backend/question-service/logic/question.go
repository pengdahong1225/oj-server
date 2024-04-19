package logic

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"question-service/logic/producer"
	pb2 "question-service/logic/proto"
	"question-service/models"
	"question-service/services/mq"
	"question-service/services/registry"
	"question-service/settings"
	"question-service/views"
	"time"
)

func QuestionSet(ctx *gin.Context, cursor int32) {
	dbConn, err := registry.NewDBConnection(settings.Conf.RegistryConfig)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "db服务连接失败",
		})
	}
	defer dbConn.Close()

	client := pb2.NewDBServiceClient(dbConn)
	request := &pb2.GetQuestionListRequest{Cursor: cursor}
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

	client := pb2.NewDBServiceClient(dbConn)
	request := &pb2.GetQuestionRequest{Id: id}
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

	client := pb2.NewDBServiceClient(dbConn)
	request := &pb2.QueryQuestionWithNameRequest{Name: name}

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

func QuestionRun(ctx *gin.Context, form *models.QuestionForm, conn *websocket.Conn) {
	// 新建上下文
	s := producer.NewSession(form.UserId)
	s.WebConnection = conn
	// 发布任务
	amqp := &producer.Amqp{
		MqConnection: mq.MqConnection,
		Exchange:     "amqp.direct",
		Queue:        "question",
		RoutingKey:   "question",
	}
	if err := amqp.Prepare(); err != nil {
		logrus.Errorf("amqp预处理失败:%s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "运行失败",
		})
		return
	}
	defer amqp.Channel.Close()
	// 序列化消息
	msg, _ := json.Marshal(form)
	if !amqp.Publish(msg) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "运行失败",
		})
		// 关闭web连接
		conn.Close()
		return
	}
	// 先返回状态
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "运行成功",
	})
	// 开启定时器
	timer := time.AfterFunc(5*time.Second, func() {
		err := conn.WriteMessage(websocket.TextMessage, []byte("执行超时"))
		if err != nil {
			logrus.Errorln("Failed to write message:", err)
		}
		conn.Close()
	})
	s.Timer = timer
}

func QuestionSubmit(ctx *gin.Context, form *models.QuestionForm, conn *websocket.Conn) {
	// 新建上下文
	s := producer.NewSession(form.UserId)
	s.WebConnection = conn
	// 发布任务
	amqp := &producer.Amqp{
		MqConnection: mq.MqConnection,
		Exchange:     "amqp.direct",
		Queue:        "question",
		RoutingKey:   "question",
	}
	if err := amqp.Prepare(); err != nil {
		logrus.Errorf("amqp预处理失败:%s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "运行失败",
		})
		return
	}
	defer amqp.Channel.Close()
	// 序列化消息
	req := &models.JudgeRequest{
		SessionID:  s.Id,
		QuestionID: form.Id,
		UserID:     form.UserId,
		Title:      form.Title,
		Code:       form.Code,
		Clang:      form.Clang,
	}
	msg, _ := json.Marshal(req)
	if !amqp.Publish(msg) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "运行失败",
		})
		// 关闭web连接
		conn.Close()
		return
	}
	// 先返回状态
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "运行成功",
	})
	// 缓存msg
	s.Msg = msg
	// 开启定时器
	timer := time.AfterFunc(5*time.Second, func() {
		err := conn.WriteMessage(websocket.TextMessage, []byte("执行超时"))
		if err != nil {
			logrus.Errorln("Failed to write message:", err)
		}
		// 记录：运行超时
		updateSubmitRecord(msg, "time out")
		conn.Close()
	})
	s.Timer = timer
}

func JudgeCallback(ctx *gin.Context, form *models.JudgeBackForm) {
	// 获取上下文
	s, ok := producer.SessionMap[form.SessionID]
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}
	// 停止定时器
	s.Timer.Stop()
	// 返回客户端执行结果
	response := views.QuestionResult{
		QuestionID: form.QuestionID,
		UserID:     form.UserID,
		Clang:      form.Clang,
		Status:     form.Status,
		Tips:       form.Tips,
		Output:     form.Output,
	}
	data, _ := json.Marshal(response)
	if err := s.WebConnection.WriteMessage(websocket.TextMessage, data); err != nil {
		logrus.Errorln("Failed to write message:", err)
	}
	// 删除上下文
	delete(producer.SessionMap, form.SessionID)
	// 关闭连接
	s.WebConnection.Close()
	// 保存提交记录
	updateSubmitRecord(s.Msg, string(data))
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

	client := pb2.NewDBServiceClient(dbConn)
	request := &pb2.UpdateUserSubmitRecordRequest{
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