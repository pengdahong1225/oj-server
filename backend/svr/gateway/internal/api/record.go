package api

import (
	"context"
	"oj-server/common/consts"
	"oj-server/common/proto/pb"
	"oj-server/svr/gateway/internal/client"
	"oj-server/svr/gateway/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RecordHandler struct {
	clients *client.Clients
}

func NewRecordHandler(clients *client.Clients) *RecordHandler {
	return &RecordHandler{
		clients: clients,
	}
}

// 处理获取用户历史提交记录
// 偏移量分页
func (rh *RecordHandler) HandleGetUserRecordList(ctx *gin.Context) {
	uid := ctx.GetInt64("uid")
	req := &model.QueryRecordListRequest{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		ResponseBadRequest(ctx, err.Error())
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	gr, err := rh.clients.Record.GetSubmitRecordList(c, &pb.GetSubmitRecordListRequest{
		Uid:      uid,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		ResponseWithGrpcError(ctx, err)
		return
	}

	ResponseOK(ctx, gr)
}

// 处理获取用户某个提交记录的具体信息
func (rh *RecordHandler) HandleGetUserRecord(ctx *gin.Context) {
	recordId, err := strconv.ParseInt(ctx.Query("record_id"), 10, 64)
	if err != nil || recordId <= 0 {
		ResponseBadRequest(ctx, "record_id invalid")
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	gr, err := rh.clients.Record.GetSubmitRecordData(c, &pb.GetSubmitRecordRequest{
		Id: recordId,
	})
	if err != nil {
		ResponseWithGrpcError(ctx, err)
		return
	}

	ResponseOK(ctx, gr)
}

// 获取排行榜
func (rh *RecordHandler) HandleGetLeaderboard(ctx *gin.Context) {

}
