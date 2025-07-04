package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"oj-server/app/gateway/internal/define"
	"oj-server/consts"
	"oj-server/module/registry"
	"oj-server/proto/pb"
	"strconv"
)

func HandleGetTagList(ctx *gin.Context) {}

func HandleGetProblemList(ctx *gin.Context) {
	response := &define.Response{
		ErrCode: pb.Error_EN_Success,
	}

	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("page_size")
	keyWord := ctx.Query("keyword")
	tag := ctx.Query("tag")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		response.ErrCode = pb.Error_EN_FormValidateFailed
		response.Message = "页码参数错误"
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		response.ErrCode = pb.Error_EN_FormValidateFailed
		response.Message = "页大小参数错误"
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	//uidStr := ctx.DefaultQuery("uid", "")
	//uid, _ := strconv.ParseInt(uidStr, 10, 64)

	// 调用problem服务
	conn, err := registry.GetGrpcConnection(consts.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		response.ErrCode = pb.Error_EN_ServiceBusy
		response.Message = "服务器错误"
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	client := pb.NewProblemServiceClient(conn)
	req := &pb.GetProblemListRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Keyword:  keyWord,
		Tag:      tag,
	}
	resp, err := client.GetProblemList(context.Background(), req)
	if err != nil {
		logrus.Errorf("problem服务获取题目列表失败:%s", err.Error())
		response.ErrCode = pb.Error_EN_Failed
		response.Message = "获取题目列表失败"
		ctx.JSON(http.StatusOK, response)
		return
	}
	response.Data = resp
	response.Message = "获取题目列表成功"
	ctx.JSON(http.StatusOK, response)
	return
}
func HandleGetProblemDetail(ctx *gin.Context) {}
func HandleSubmitProblem(ctx *gin.Context)    {}
func HandleGetSubmitResult(ctx *gin.Context)  {}
func HandleUpdateProblem(ctx *gin.Context)    {}
func HandleCreateProblem(ctx *gin.Context)    {}
func HandleDeleteProblem(ctx *gin.Context)    {}
