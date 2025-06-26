package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pengdahong1225/oj-server/backend/app/gateway/internal/define"
	"github.com/pengdahong1225/oj-server/backend/consts"
	"github.com/pengdahong1225/oj-server/backend/module/registry"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func HandleGetTagList(ctx *gin.Context) {}

func HandleGetProblemList(ctx *gin.Context) {
	response := &define.Response{}
	defer ctx.JSON(http.StatusOK, response)

	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("page_size")
	keyWord := ctx.Query("keyword")
	tag := ctx.Query("tag")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		response.Code = define.Failed
		response.Message = "页码参数错误"
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		response.Code = define.Failed
		response.Message = "页大小参数错误"
		return
	}
	//uidStr := ctx.DefaultQuery("uid", "")
	//uid, _ := strconv.ParseInt(uidStr, 10, 64)

	// 调用problem服务
	conn, err := registry.GetGrpcConnection(consts.ProblemService)
	if err != nil {
		logrus.Errorf("problem服务连接失败:%s", err.Error())
		response.Code = define.Failed
		response.Message = "problem服务连接失败"
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
		response.Code = define.Failed
		response.Message = "problem服务获取题目列表失败"
		return
	}
	response.Code = define.Success
	response.Data = resp
	return
}
func HandleGetProblemDetail(ctx *gin.Context) {}
func HandleSubmitProblem(ctx *gin.Context)    {}
func HandleGetSubmitResult(ctx *gin.Context)  {}
func HandleUpdateProblem(ctx *gin.Context)    {}
func HandleCreateProblem(ctx *gin.Context)    {}
func HandleDeleteProblem(ctx *gin.Context)    {}
