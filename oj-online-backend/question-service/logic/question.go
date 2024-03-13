package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"question-service/global"
	"question-service/models"
	pb "question-service/proto"
	"time"
)

var ConnectionMap = make(map[string]*websocket.Conn, 10) // sessionID,websocket连接
// 保存上下文
func loadContext(conn *websocket.Conn, form *models.QuestionForm) {
	sessionID := fmt.Sprintf("%d-%d", form.UserId, time.Now().Unix())
	ConnectionMap[sessionID] = conn
}

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
	// 保存上下文
	loadContext(conn, form)
	// 发布任务
	amqp := &Amqp{
		MqConnection: global.MqConnection,
		exchange:     "amqp.direct",
		queue:        "question",
		routingKey:   "question",
	}
	if amqp.prepare(ctx) {
		// 序列化消息
		msg, _ := json.Marshal(form)
		if amqp.publish(msg) {
			ctx.JSON(http.StatusOK, gin.H{
				"msg": "运行成功",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "运行失败",
			})
		}
		// 关闭通道
		amqp.channel.Close()
	}
}

func QuestionSubmit(ctx *gin.Context, form *models.QuestionForm) {

}

func JudgeCallback(ctx *gin.Context, sessionID string) {
	// 获取上下文
	conn, _ := ConnectionMap[sessionID]
	defer conn.Close()
	// 返回客户端执行结果

}
