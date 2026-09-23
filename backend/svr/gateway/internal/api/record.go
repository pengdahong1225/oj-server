package api

import (
	"context"
	"oj-server/common/consts"
	"oj-server/common/proto/pb"
	"oj-server/svr/gateway/internal/client"
	"oj-server/svr/gateway/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RecordHandler struct {
	clients *client.Clients
}

func NewRecordHandler(clients *client.Clients) *RecordHandler {
	return &RecordHandler{
		clients: clients,
	}
}

// 获取判题任务结果
// 拿到判题任务结果，再去获取任务的output
func (rh *RecordHandler) HandleGetSubmitResult(ctx *gin.Context) {
	taskId, exists := ctx.GetQuery("task_id")
	if !exists || taskId == "" {
		ResponseBadRequest(ctx, "task_id 不能为空")
		return
	}
	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	resp, err := rh.clients.Record.QueryJudgeResult(c, &pb.QueryJudgeResultRequest{})
	if err != nil {
		logrus.Errorf("获取判题结果失败:%s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	result := &model.JudgeResultAbstract{
		Accepted: resp.Accepted,
		Message:  resp.Message,
	}
	ResponseOK(ctx, result)
}

// 处理获取用户历史提交记录
// 偏移量分页
func (rh *RecordHandler) HandleGetUserRecordList(ctx *gin.Context) {
	uid := ctx.GetInt64("uid")
	var params model.QueryUserRecordListParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ResponseBadRequest(ctx, "参数验证失败")
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	resp, err := rh.clients.Record.GetSubmitRecordList(c, &pb.GetSubmitRecordListRequest{
		Uid:      uid,
		Page:     params.Page,
		PageSize: params.PageSize,
	})
	if err != nil {
		logrus.Errorf("获取提交记录失败:%s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}
	result := &model.QueryUserRecordListResult{
		Total: int64(resp.Total),
	}
	for _, pbRecord := range resp.Data {
		var record model.Record
		record.FromPbRecord(pbRecord)
		result.List = append(result.List, record)
	}

	ResponseOK(ctx, result)
}

// 处理获取用户某个提交记录的具体信息
func (rh *RecordHandler) HandleGetUserRecord(ctx *gin.Context) {
	recordId, err := strconv.ParseInt(ctx.Query("record_id"), 10, 64)
	if err != nil || recordId <= 0 {
		ResponseBadRequest(ctx, "record_id 不能为空")
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), consts.GrpcTimeout)
	defer cancel()
	resp, err := rh.clients.Record.GetSubmitRecordData(c, &pb.GetSubmitRecordRequest{Id: recordId})
	if err != nil {
		logrus.Errorf("查询提交记录失败:%s", err.Error())
		ResponseWithGrpcError(ctx, err)
		return
	}

	record := &model.Record{}
	record.FromPbRecord(resp.Data)
	ResponseOK(ctx, record)
}

// 处理排行榜查询
func (rh *RecordHandler) HandleGetLeaderboard(ctx *gin.Context) {

}

// 处理获取用户的AC题目列表
func (rh *RecordHandler) HandleGetUserSolvedList(ctx *gin.Context) {}
