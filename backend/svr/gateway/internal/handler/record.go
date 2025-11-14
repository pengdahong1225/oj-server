package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"oj-server/global"
	"oj-server/pkg/proto/pb"
	"oj-server/pkg/registry"
	"oj-server/svr/gateway/internal/model"
	"strconv"
)

// 处理排行榜查询
func HandleGetLeaderboard(ctx *gin.Context) {

}

// 处理获取用户的AC题目列表
func HandleGetUserSolvedList(ctx *gin.Context) {}

// 获取判题任务结果
// 拿到判题任务结果，再去获取任务的output
func HandleGetSubmitResult(ctx *gin.Context) {}

// 处理获取用户历史提交记录
// 偏移量分页
func HandleGetUserRecordList(ctx *gin.Context) {
	// 获取元数据
	uid := ctx.GetInt64("uid")

	// 查询参数校验
	var params model.QueryUserRecordListParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ResponseWithJson(ctx, http.StatusBadRequest, "参数验证失败", nil)
		return
	}

	// 调用题目服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewRecordServiceClient(conn)
	request := &pb.GetSubmitRecordListRequest{
		Uid:      uid,
		Page:     params.Page,
		PageSize: params.PageSize,
	}
	resp, err := client.GetSubmitRecordList(context.Background(), request)
	if err != nil {
		logrus.Errorf("获取提交记录失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "获取提交记录失败", nil)
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

	ResponseWithJson(ctx, http.StatusOK, "success", result)
}

// 处理获取用户某个提交记录的具体信息
func HandleGetUserRecord(ctx *gin.Context) {
	recordId, err := strconv.ParseInt(ctx.Query("record_id"), 10, 64)
	if err != nil || recordId <= 0 {
		ResponseWithJson(ctx, http.StatusBadRequest, "record_id 不能为空", nil)
		return
	}
	// 调用题目服务
	conn, err := registry.MyRegistrar.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "服务繁忙", nil)
		return
	}
	client := pb.NewRecordServiceClient(conn)
	rpc_resp, err := client.GetSubmitRecordData(ctx, &pb.GetSubmitRecordRequest{
		Id: recordId,
	})
	if err != nil {
		logrus.Errorf("查询提交记录失败:%s", err.Error())
		ResponseWithJson(ctx, http.StatusInternalServerError, "查询提交记录失败", nil)
		return
	}

	record := &model.Record{}
	record.FromPbRecord(rpc_resp.Data)
	ResponseWithJson(ctx, http.StatusOK, "success", record)
}
