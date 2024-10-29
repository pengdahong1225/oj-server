package logic

import (
	"context"
	"github.com/pengdahong1225/oj-server/backend/app/question-service/internal/models"
	"github.com/pengdahong1225/oj-server/backend/module/registry"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
)

type NoticeLogic struct {
}

func (r NoticeLogic) GetNoticeList(params *models.QueryNoticeListParams) *models.Response {
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

	client := pb.NewNoticeServiceClient(dbConn)
	response, err := client.GetNoticeList(context.Background(), &pb.GetNoticeListRequest{
		Page:     params.Page,
		PageSize: params.PageSize,
		Keyword:  params.KeyWord,
	})
	if err != nil {
		res.Code = models.Failed
		res.Message = err.Error()
		return res
	}
	data := &models.NoticeRspData{
		Total:      response.Total,
		NoticeList: response.Data,
	}

	res.Data = data
	res.Message = "OK"
	return res
}
