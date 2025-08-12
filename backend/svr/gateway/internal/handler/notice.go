package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"oj-server/proto/pb"
	"oj-server/src/gateway/internal/define"
)

func HandleGetNoticeList(ctx *gin.Context) {
	resp := define.Response{
		ErrCode: pb.Error_EN_Success,
		Message: "",
		Data:    nil,
	}

	// 表单验证
	from := &define.QueryNoticeListParams{}
	err := ctx.ShouldBindQuery(from)
	if err != nil {
		logrus.Debugf("表单验证失败: %v", err)
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "表单验证失败"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// ...

	resp.Data = &define.QueryNoticeListResponse{
		NoticeList: nil,
		Total:      0,
	}

	ctx.JSON(http.StatusOK, resp)
}
