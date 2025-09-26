package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"oj-server/global"
	"oj-server/module/registry"
	"oj-server/proto/pb"
	"oj-server/svr/gateway/internal/model"
	"strconv"
)

// 处理排行榜查询
func HandleGetLeaderboard(ctx *gin.Context) {
	resp := &model.Response{
		ErrCode: pb.Error_EN_Success,
		Message: "success",
	}
}

// 处理获取用户的AC题目列表
func HandleGetUserSolvedList(ctx *gin.Context) {}

// 获取判题任务结果
// 拿到判题任务结果，再去获取任务的output
func HandleGetSubmitResult(ctx *gin.Context) {}

// 处理获取用户历史提交记录
// 偏移量分页
func HandleGetUserRecordList(ctx *gin.Context) {
	resp := &model.Response{
		ErrCode: pb.Error_EN_Success,
		Message: "success",
	}

	// 获取元数据
	uid := ctx.GetInt64("uid")

	// 查询参数校验
	var params model.QueryUserRecordListParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "参数验证失败"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	conn, err := registry.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务繁忙"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	client := pb.NewRecordServiceClient(conn)
	request := &pb.GetSubmitRecordListRequest{
		Uid:      uid,
		Page:     params.Page,
		PageSize: params.PageSize,
	}
	rpc_resp, err := client.GetSubmitRecordList(context.Background(), request)
	if err != nil {
		logrus.Errorf("获取提交记录失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = "获取提交记录失败"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	data := &model.QueryUserRecordListResult{
		Total: int64(rpc_resp.Total),
	}
	for _, pbRecord := range rpc_resp.Data {
		var record model.Record
		record.FromPbRecord(pbRecord)
		data.List = append(data.List, record)
	}
	resp.Data = data
	ctx.JSON(http.StatusOK, resp)
}

// 处理获取用户某个提交记录的具体信息
func HandleGetUserRecord(ctx *gin.Context) {
	resp := &model.Response{
		ErrCode: pb.Error_EN_Success,
		Message: "success",
	}

	recordId, err := strconv.ParseInt(ctx.Query("record_id"), 10, 64)
	if err != nil || recordId <= 0 {
		resp.ErrCode = pb.Error_EN_FormValidateFailed
		resp.Message = "record_id 不能为空"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	// 调用题目服务
	conn, err := registry.GetGrpcConnection(global.ProblemService)
	if err != nil {
		logrus.Errorf("用户服务连接失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_ServiceBusy
		resp.Message = "服务繁忙"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	client := pb.NewRecordServiceClient(conn)
	rpc_resp, err := client.GetSubmitRecordData(ctx, &pb.GetSubmitRecordRequest{
		Id: recordId,
	})
	if err != nil {
		logrus.Errorf("查询提交记录失败:%s", err.Error())
		resp.ErrCode = pb.Error_EN_Failed
		resp.Message = "查询提交记录失败"
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	record := &model.Record{}
	record.FromPbRecord(rpc_resp.Data)
	resp.Data = record
	ctx.JSON(http.StatusOK, resp)
}
