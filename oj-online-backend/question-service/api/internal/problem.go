package internal

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	pb "question-service/api/proto"
	"question-service/models"
	"question-service/services/registry"
	"question-service/settings"
)

func ProblemSet(cursor int) *models.Response {
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

	// 序列化
	marshalOptions := protojson.MarshalOptions{
		UseProtoNames:     true,
		UseEnumNumbers:    true,
		EmitDefaultValues: true,
	}
	bytes, err := marshalOptions.Marshal(response)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		return res
	}

	res.Code = http.StatusOK
	res.Message = "OK"
	res.Data = string(bytes)
	return res
}

func QuestionDetail(id int64) {

}

func QuestionQuery(name string) {

}

func QuestionRun(form *models.QuestionForm) {

}
func QuestionSubmit(ctx *gin.Context, form *models.QuestionForm) {

}

// 保存提交记录
func updateSubmitRecord(msg []byte, result string) {

}
