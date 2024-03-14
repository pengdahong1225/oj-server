package logic

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"question-service/global"
	"question-service/logic/internal"
	"question-service/models"
	pb "question-service/proto"
	"time"
)

func QuestionSet(ctx *gin.Context, cursor int32) {
	dbConn, err := global.NewDBConnection()
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
	var questionList []models.Question
	for _, v := range response.Data {
		questionList = append(questionList, models.Question{
			Id:          v.Id,
			Title:       v.Title,
			Description: v.Description,
			Tags:        v.Tags,
			Level:       v.Level,
			CreateAt:    v.CreateAt,
		})
	}
	data, _ := json.Marshal(questionList)
	ctx.JSON(http.StatusOK, gin.H{
		"data":   data,
		"cursor": response.Cursor,
	})
}

func QuestionDetail(ctx *gin.Context, id int64) {
	dbConn, err := global.NewDBConnection()
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
	data := models.Question{
		Id:          response.Data.Id,
		Title:       response.Data.Title,
		Description: response.Data.Description,
		Tags:        response.Data.Tags,
		Level:       response.Data.Level,
		CreateAt:    response.Data.CreateAt,
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func QuestionQuery(ctx *gin.Context, name string) {
	dbConn, err := global.NewDBConnection()
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
	var questionList []models.Question
	for _, v := range response.Data {
		questionList = append(questionList, models.Question{
			Id:          v.Id,
			Title:       v.Title,
			Description: v.Description,
			Tags:        v.Tags,
			Level:       v.Level,
			CreateAt:    v.CreateAt,
		})
	}
	data, _ := json.Marshal(questionList)
	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func QuestionRun(ctx *gin.Context, form *models.QuestionForm, conn *websocket.Conn) {
	// 新建上下文
	s := internal.NewSession(form.UserId)
	s.WebConnection = conn
	// 发布任务
	amqp := &internal.Amqp{
		MqConnection: global.MqConnection,
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
	s := internal.NewSession(form.UserId)
	s.WebConnection = conn
	// 发布任务
	amqp := &internal.Amqp{
		MqConnection: global.MqConnection,
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
	s, ok := internal.SessionMap[form.SessionID]
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}
	// 停止定时器
	s.Timer.Stop()
	// 返回客户端执行结果
	response := models.QuestionResult{
		Id:     form.Id,
		UserId: form.UserId,
		Clang:  form.Clang,
		Result: form.Result,
	}
	data, _ := json.Marshal(response)
	if err := s.WebConnection.WriteMessage(websocket.TextMessage, data); err != nil {
		logrus.Errorln("Failed to write message:", err)
	}
	// 删除上下文
	delete(internal.SessionMap, form.SessionID)
	// 关闭连接
	s.WebConnection.Close()
	// 保存提交记录
	updateSubmitRecord(s.Msg, form.Result)
}

// 保存提交记录
func updateSubmitRecord(msg []byte, result string) {
	dbConn, err := global.NewDBConnection()
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
