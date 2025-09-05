package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"oj-server/module/proto/pb"
	"oj-server/svr/gateway/internal/model"
)

func HandleGetNoticeList(ctx *gin.Context) {
	resp := model.Response{
		ErrCode: pb.Error_EN_Success,
		Message: "",
		Data:    nil,
	}

	// 参数验证
	from := &model.QueryNoticeListParams{}
	err := ctx.ShouldBindQuery(from)
	if err != nil {
		logrus.Debugf("表单验证失败: %v", err)
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "表单验证失败"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// ...

	resp.Data = &model.QueryNoticeListResponse{
		NoticeList: nil,
		Total:      0,
	}

	ctx.JSON(http.StatusOK, resp)
}
